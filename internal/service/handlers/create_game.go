package handlers

import (
	"encoding/json"
	"github.com/olegfomenko/game-service/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/go/xdrbuild"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// TODO change constant
const IssuanceAmount = uint64(1000000000000000)

func CreateGame(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewCreateGame(r)
	if err != nil {
		Log(r).WithError(err).Error("invalid request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	details, err := json.Marshal(map[string]interface{}{
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

	date := strings.Split(request.Data.Attributes.Date, " ")
	date[0] = strings.ReplaceAll(date[0], "-", "")
	assetCode := "GAM" + date[0]

	amount := uint64(request.Data.Attributes.Amount)

	assetType, err := Connector(r).GetUint32KeyValue("asset_type:gam")
	if err != nil {
		Log(r).WithError(err).Error("error getting key value for asset_type:gam")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	Log(r).Info("Got gam asset_type=", assetType)

	// TODO check TrailingDigitsCount
	createGam := &xdrbuild.CreateAsset{
		RequestID:                0,
		Code:                     assetCode,
		MaxIssuanceAmount:        IssuanceAmount,
		PreIssuanceSigner:        Connector(r).Signer().Address(),
		InitialPreIssuanceAmount: IssuanceAmount,
		TrailingDigitsCount:      6,
		Policies:                 0,
		Type:                     uint64(assetType),
		CreatorDetails:           json.RawMessage(details),
		AllTasks:                 nil,
	}

	createGamBalance := &xdrbuild.ManageBalanceOp{
		Action:      xdr.ManageBalanceActionCreate,
		Destination: request.Data.Attributes.OwnerId,
		AssetCode:   assetCode,
	}

	respTx, err := Connector(r).SubmitSigned(r.Context(), createGam, createGamBalance)
	if err != nil {
		Log(r).WithError(err).Error("error sending transaction")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	gamBalance, err := Connector(r).Balance(request.Data.Attributes.OwnerId, assetCode)
	if err != nil {
		Log(r).WithError(err).Error("error getting organizer gam balance")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	Log(r).Info("Got organizer gam balance", gamBalance.ID)

	// TODO fix  reference
	issueGam := &xdrbuild.CreateIssuanceRequest{
		Reference: strconv.Itoa(time.Now().Nanosecond()),
		Receiver:  request.Data.Attributes.OwnerId,
		Asset:     gamBalance.ID,
		Amount:    amount,
		Details:   json.RawMessage(details),
		AllTasks:  nil,
	}
	Log(r).Info("Issuing asset operation:", issueGam)

	orgBalance, err := Connector(r).GetBalanceByID(request.Data.Attributes.SourceBalanceId)
	if err != nil {
		Log(r).WithError(err).Error("error getting organizer usd balance")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	Log(r).Info("Got organizer balance", orgBalance.ID)

	// TODo fix reference
	paymentToAdmin, _ := xdrbuild.CreatePaymentForAccount(xdrbuild.CreatePaymentForAccountOpts{
		SourceAccountID:      &orgBalance.Relationships.Owner.Data.ID,
		SourceBalanceID:      orgBalance.ID,
		DestinationAccountID: Connector(r).Signer().Address(),
		Amount:               amount,
		Subject:              "Creating game organizer payment",
		Reference:            strconv.Itoa(time.Now().Nanosecond()),
		Fee: xdrbuild.Fee{
			SourceFixed:        0,
			SourcePercent:      0,
			DestinationFixed:   0,
			DestinationPercent: 0,
			SourcePaysForDest:  false,
		},
	})
	Log(r).Info("Payment for asset operation:", paymentToAdmin)

	respTx, err = Connector(r).SubmitSigned(r.Context(), issueGam, paymentToAdmin)
	if err != nil {
		Log(r).WithError(err).Error("error sending transaction")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	ape.Render(w, respTx)
}
