package handlers

import (
	"encoding/json"
	"github.com/olegfomenko/game-service/internal/horizon"
	"github.com/olegfomenko/game-service/internal/service/requests"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/go/xdrbuild"
	regources "gitlab.com/tokend/regources/generated"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const IssuanceAmount = uint64(1000_000_000_000_000)

func CreateGame(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewCreateGame(r)
	if err != nil {
		Log(r).WithError(err).Error("invalid request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	details, err := parseDetails(map[string]interface{}{
		"organizer":        request.Data.Attributes.OwnerId,
		"date":             request.Data.Attributes.Date,
		"team1":            request.Data.Attributes.Team1,
		"team2":            request.Data.Attributes.Team2,
		"stream_link":      request.Data.Attributes.StreamLink,
		"name_competition": request.Data.Attributes.NameCompetition,
	})
	if err != nil {
		Log(r).WithError(err).Error("marshaling details to JSON")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	date, err := parseGameDate(request.Data.Attributes.Date)
	if err != nil {
		Log(r).WithError(err).Error("error parsing date")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	assetCode := "GAM" + date
	Log(r).Info("Creating asset ", assetCode)

	assetType, err := Connector(r).GetUint32KeyValue("asset_type:gam")
	if err != nil {
		Log(r).WithError(err).Error("error getting key value for asset_type:gam")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	Log(r).Info("Got GAM asset_type=", assetType)

	createGam := &xdrbuild.CreateAsset{
		RequestID:                0,
		Code:                     assetCode,
		MaxIssuanceAmount:        IssuanceAmount,
		PreIssuanceSigner:        Connector(r).Source().Address(),
		InitialPreIssuanceAmount: IssuanceAmount,
		TrailingDigitsCount:      6,
		Policies:                 0,
		Type:                     uint64(assetType),
		CreatorDetails:           details,
		AllTasks:                 nil,
	}

	respTx, err := Connector(r).SubmitSigned(r.Context(), nil, createGam)
	if err != nil {
		Log(r).WithError(err).Error("error sending transaction")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	respTx, err = donate(
		r,
		request.Data.Attributes.OwnerId,
		assetCode,
		uint64(request.Data.Attributes.Amount),
		request.Data.Attributes.SourceBalanceId,
		details,
	)

	if err != nil {
		Log(r).WithError(err).Error("error donating")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	ape.Render(w, respTx)
}

func loadBalance(r *http.Request, assetCode string, ownerID string) (*regources.Balance, error) {
	gamBalance, err := Connector(r).Balance(ownerID, assetCode)

	if err == horizon.ErrNotFound || err == horizon.ErrDataEmpty {
		Log(r).Info("Crating new balance for user ", ownerID)
		createBalance := &xdrbuild.ManageBalanceOp{
			Action:      xdr.ManageBalanceActionCreate,
			Destination: ownerID,
			AssetCode:   assetCode,
		}

		_, err := Connector(r).SubmitSigned(r.Context(), nil, createBalance)
		if err != nil {
			return nil, err
		}

		gamBalance, err := Connector(r).Balance(ownerID, assetCode)
		if err != nil {
			return nil, err
		}

		return gamBalance, nil
	} else if err != nil {
		return nil, err
	}

	return gamBalance, nil
}

func donate(r *http.Request, ownerID string, gamAsset string, amount uint64, sourceBalance string, details json.RawMessage) (*regources.TransactionResponse, error) {
	gamBalance, err := loadBalance(r, gamAsset, ownerID)
	if err != nil {
		return nil, err
	}
	Log(r).Info("Got organizer gam balance", gamBalance.ID)

	tasks := uint32(1)

	paymentReference := strconv.Itoa(time.Now().Nanosecond())
	payment := xdrbuild.CreateRedemptionRequest{
		SourceBalanceID:      sourceBalance,
		DestinationAccountID: Connector(r).Signer().Address(),
		Amount:               amount,
		Reference:            paymentReference,
		Details:              details,
		AllTasks:             &tasks,
	}
	Log(r).Info("Payment operation:", payment)

	respTx, err := Connector(r).SubmitSigned(r.Context(), nil, payment)
	if err != nil {
		return nil, err
	}

	issueGam := &xdrbuild.CreateIssuanceRequest{
		Reference: strconv.Itoa(time.Now().Nanosecond()),
		Receiver:  gamBalance.ID,
		Asset:     gamAsset,
		Amount:    amount,
		Details:   details,
		AllTasks:  nil,
	}
	Log(r).Info("Issuing asset operation:", issueGam)

	requests, err := Connector(r).GetRedemptionRequests(21, 1)
	if err != nil {
		return nil, err
	}

	var req *regources.ReviewableRequest = nil
	for _, r := range requests {
		if *r.Attributes.Reference == paymentReference {
			req = &r
		}
	}
	if req == nil {
		return nil, errors.New("cannot find redemption")
	}

	id, _ := strconv.Atoi(req.ID)
	reviewOrg := xdrbuild.ReviewRequest{
		ID:      uint64(id),
		Hash:    &req.Attributes.Hash,
		Action:  xdr.ReviewRequestOpActionApprove,
		Details: xdrbuild.ReviewableRequestBaseDetails{RequestType: xdr.ReviewableRequestTypePerformRedemption},
		ReviewDetails: xdrbuild.ReviewDetails{
			TasksToAdd:      0,
			TasksToRemove:   tasks,
			ExternalDetails: "",
		},
	}

	Log(r).Info("Review operation:", reviewOrg)
	Log(r).Info(id, " ", req.Attributes.Hash, " ", *req.Attributes.Reference)

	respTx, err = Connector(r).SubmitSigned(r.Context(), nil, reviewOrg, issueGam)
	if err != nil {
		return nil, err
	}

	return respTx, nil
}

func parseGameDate(date string) (string, error) {
	arr := strings.Split(date, " ")
	if len(arr) == 0 {
		return "", errors.New("invalid date string")
	}

	arr[0] = strings.ReplaceAll(arr[0], "-", "")
	return arr[0], nil
}

func parseDetails(details map[string]interface{}) (json.RawMessage, error) {
	data, err := json.Marshal(details)

	if err != nil {
		return nil, err
	}

	return data, nil
}
