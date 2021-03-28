package handlers

import (
	"encoding/json"
	"github.com/olegfomenko/game-service/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/go/xdrbuild"
	regources "gitlab.com/tokend/regources/generated"
	"net/http"
	"strconv"
	"time"
)

func PayGame(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewPayGame(r)
	if err != nil {
		Log(r).WithError(err).Error("invalid request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	assetCode := request.Data.Attributes.GameCoinId
	amount := uint64(request.Data.Attributes.Amount)

	// CREATING GAM BALANCE
	createGamBalance := &xdrbuild.ManageBalanceOp{
		Action:      xdr.ManageBalanceActionCreate,
		Destination: request.Data.Attributes.OwnerId,
		AssetCode:   assetCode,
	}

	respTx, err := Connector(r).SubmitSigned(r.Context(), nil, createGamBalance)
	if err != nil {
		Log(r).WithError(err).Error("error sending transaction")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	// GETTING CREATED GAM BALANCE
	gamBalance, err := Connector(r).Balance(request.Data.Attributes.OwnerId, assetCode)
	if err != nil {
		Log(r).WithError(err).Error("error getting organizer gam balance")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	Log(r).Info("Got organizer gam balance", gamBalance.ID)

	// CREATING PAYMENT USER->ADMIN
	tasks := uint32(1)
	paymentReference := strconv.Itoa(time.Now().Nanosecond())
	payment := xdrbuild.CreateRedemptionRequest{
		SourceBalanceID:      request.Data.Attributes.SourceBalanceId,
		DestinationAccountID: Connector(r).Signer().Address(),
		Amount:               amount,
		Reference:            paymentReference,
		Details:              json.RawMessage{},
		AllTasks:             &tasks,
	}
	Log(r).Info("Payment operation:", payment)

	respTx, err = Connector(r).SubmitSigned(r.Context(), nil, payment)
	if err != nil {
		Log(r).WithError(err).Error("error sending transaction")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	// CREATING ISSUE GAM OP
	issueGam := &xdrbuild.CreateIssuanceRequest{
		Reference: strconv.Itoa(time.Now().Nanosecond()),
		Receiver:  gamBalance.ID,
		Asset:     assetCode,
		Amount:    amount,
		Details:   json.RawMessage{},
		AllTasks:  nil,
	}
	Log(r).Info("Issuing asset operation:", issueGam)

	// GETTING PENDING REDEMPTIONS
	requests, err := Connector(r).GetRedemptionRequests(21, 1)
	var req *regources.ReviewableRequest = nil

	for _, r := range requests {
		if *r.Attributes.Reference == paymentReference {
			req = &r
		}
	}

	id, _ := strconv.Atoi(req.ID)

	// CREATING APPROVE OP
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

	// SUBMITTING APPROVE & ISSUANCE OPs
	respTx, err = Connector(r).SubmitSigned(r.Context(), nil, reviewOrg, issueGam)
	if err != nil {
		Log(r).WithError(err).Error("error sending transaction")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	ape.Render(w, respTx)
}
