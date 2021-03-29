package handlers

import (
	"encoding/json"
	"github.com/pkg/errors"
	"gitlab.com/tokend/go/xdr"
	"gitlab.com/tokend/go/xdrbuild"
	regources "gitlab.com/tokend/regources/generated"
	"net/http"
	"strconv"
	"time"
)

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

func parseDetails(details map[string]interface{}) (json.RawMessage, error) {
	data, err := json.Marshal(details)

	if err != nil {
		return nil, err
	}

	return data, nil
}
