package handlers

import (
	"github.com/olegfomenko/game-service/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/go/xdrbuild"
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

	//TODO maybe use details as comment with donate
	//TODO fix reference
	issueGam := &xdrbuild.CreateIssuanceRequest{
		Reference: strconv.Itoa(time.Now().Nanosecond()),
		Receiver:  request.Data.Attributes.OwnerId,
		Asset:     assetCode,
		Amount:    amount,
		Details:   nil,
		AllTasks:  nil,
	}

	orgBalance, err := Connector(r).GetBalanceByID(request.Data.Attributes.SourceBalanceId)
	if err != nil {
		Log(r).WithError(err).Error("error getting organizer balance")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	// TODO fix reference
	paymentToAdmin, _ := xdrbuild.CreatePaymentForAccount(xdrbuild.CreatePaymentForAccountOpts{
		SourceAccountID:      &orgBalance.Relationships.Owner.Data.ID,
		SourceBalanceID:      orgBalance.ID,
		DestinationAccountID: Connector(r).Signer().Address(),
		Amount:               amount,
		Subject:              "Game payment from user",
		Reference:            strconv.Itoa(time.Now().Nanosecond()),
		Fee: xdrbuild.Fee{
			SourceFixed:        0,
			SourcePercent:      0,
			DestinationFixed:   0,
			DestinationPercent: 0,
			SourcePaysForDest:  false,
		},
	})

	_, err = Connector(r).TxSigned(Connector(r).Signer(), issueGam, paymentToAdmin)
	if err != nil {
		Log(r).WithError(err).Error("error sending transaction")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	w.WriteHeader(http.StatusOK)
}
