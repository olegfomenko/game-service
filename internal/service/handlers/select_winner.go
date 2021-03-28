package handlers

import (
	"encoding/json"
	"github.com/olegfomenko/game-service/internal/service/requests"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/go/xdrbuild"
	regources "gitlab.com/tokend/regources/generated"
	"net/http"
	"strconv"
	"time"
)

func SelectWinner(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewSelectWinner(r)
	if err != nil {
		Log(r).WithError(err).Error("invalid request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	gam, err := Connector(r).Asset(request.Data.Attributes.GameCoinId)
	if err != nil {
		Log(r).WithError(err).Error("error getting GAm asset")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	amount := uint64(gam.Attributes.Issued)
	amountPerPlayer := amount / 5

	var details = make(map[string]interface{})

	err = json.Unmarshal(gam.Attributes.Details, &details)
	if err != nil {
		Log(r).WithError(err).Error("error while parsing details")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	team1 := details["team1"].(map[string]string)
	team2 := details["team2"].(map[string]string)

	if team1["name"] == request.Data.Attributes.TeamName {
		_, err := payPrizes(r, amountPerPlayer, team1)

		if err != nil {
			Log(r).WithError(err).Error("error submitting tx")
			ape.RenderErr(w, problems.BadRequest(err)...)
			return
		}
	} else {
		_, err := payPrizes(r, amountPerPlayer, team2)

		if err != nil {
			Log(r).WithError(err).Error("error submitting tx")
			ape.RenderErr(w, problems.BadRequest(err)...)
			return
		}
	}

	details["winner"] = request.Data.Attributes.TeamName
	raw, _ := parseDetails(details)

	updateAsset := &xdrbuild.UpdateAsset{
		Code:           request.Data.Attributes.GameCoinId,
		CreatorDetails: raw,
		AllTasks:       nil,
	}

	respTx, err := Connector(r).SubmitSigned(r.Context(), nil, updateAsset)
	if err != nil {
		Log(r).WithError(err).Error("error updating asset")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	ape.Render(w, respTx)
}

func payPrizes(r *http.Request, amount uint64, team map[string]string) (*regources.TransactionResponse, error) {
	var payments []xdrbuild.Operation

	source := Connector(r).Source().Address()
	signer := Connector(r).Signer().Seed()

	for k, id := range team {
		if k != "name" {
			payment, err := xdrbuild.CreatePaymentForAccount(xdrbuild.CreatePaymentForAccountOpts{
				SourceAccountID:      &source,
				SourceBalanceID:      signer,
				DestinationAccountID: id,
				Amount:               amount,
				Subject:              "Game winner prize",
				Reference:            strconv.Itoa(time.Now().Nanosecond()),
				Fee: xdrbuild.Fee{
					SourceFixed:        0,
					SourcePercent:      0,
					DestinationFixed:   0,
					DestinationPercent: 0,
					SourcePaysForDest:  false,
				},
			})

			if err != nil {
				return nil, errors.Wrap(err, "error creating payment op")
			}

			payments = append(payments, payment)
		}
	}

	return Connector(r).SubmitSigned(r.Context(), nil, payments...)
}
