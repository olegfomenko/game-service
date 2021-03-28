package handlers

import (
	"encoding/json"
	"github.com/olegfomenko/game-service/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func PayTeam(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewPayTeam(r)
	if err != nil {
		Log(r).WithError(err).Error("invalid request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	dateSuffix := request.Data.Attributes.GameCoinId[3:]

	tam1ID := "TAM1" + dateSuffix
	tam2ID := "TAM2" + dateSuffix

	gam, err := Connector(r).Asset(request.Data.Attributes.GameCoinId)
	if err != nil {
		Log(r).WithError(err).Error("error getting gam")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	var details = make(map[string]interface{})
	json.Unmarshal(gam.Attributes.Details, &details)

	team1 := details["team1"].(map[string]string)
	team2 := details["team2"].(map[string]string)

	if team1["name"] == request.Data.Attributes.TeamName {

	}

	respTx, err := donate(
		r,
		request.Data.Attributes.OwnerId,
		tam1ID,
		uint64(request.Data.Attributes.Amount),
		request.Data.Attributes.SourceBalanceId,
		json.RawMessage(gam.Attributes.Details),
	)
	if err != nil {
		Log(r).WithError(err).Error("error donating")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	ape.Render(w, respTx)
}
