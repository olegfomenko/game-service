package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/olegfomenko/game-service/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/go/xdrbuild"
	"net/http"
	"strconv"
	"time"
)

// TODO op1: Pay TeamGam admin -> user, op2: Pay usd user -> admin
func PayTeam(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewPayTeam(r)
	if err != nil {
		Log(r).WithError(err).Error("invalid request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	Gam := request.Data.Attributes.GameAssetId
	asset,err := Connector(r).Asset(Gam)
	var details map[string]interface{}
	if err := json.Unmarshal(asset.Attributes.Details, &details); err != nil {
		Log(r).WithError(err).Error("cannot unmarshal json")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	assetId1 := fmt.Sprintf("%v", details["tam1"])
	assetId2 := fmt.Sprintf("%v", details["tam2"])
	team1 := fmt.Sprintf("%v", details["team1"])
	team2 := fmt.Sprintf("%v", details["team2"])
	if assetId1 == "" || assetId2 == "" || team1=="" || team2 == "" {
		Log(r).WithError(err).Error("cannot get team or tam from details")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	var finalAssetId string
	switch request.Data.Attributes.TeamName {
	case team1:
		finalAssetId = assetId1
	case team2:
		finalAssetId = assetId2
	default:
		Log(r).WithError(err).Error("incorrect team found in the gam")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	amount := uint64(request.Data.Attributes.Amount)
	//TODO maybe use details as comment with donate
	//TODO fix reference
	issueTam := &xdrbuild.CreateIssuanceRequest{
		Reference: strconv.Itoa(time.Now().Nanosecond()),
		Receiver:  request.Data.Attributes.OwnerId,
		Asset:     finalAssetId,
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
		Subject:              "Team payment from user",
		Reference:            strconv.Itoa(time.Now().Nanosecond()),
		Fee: xdrbuild.Fee{
			SourceFixed:        0,
			SourcePercent:      0,
			DestinationFixed:   0,
			DestinationPercent: 0,
			SourcePaysForDest:  false,
		},
	})

	_, err = Connector(r).TxSigned(Connector(r).Signer(), issueTam, paymentToAdmin)
	if err != nil {
		Log(r).WithError(err).Error("error sending transaction")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	w.WriteHeader(http.StatusOK)
}
