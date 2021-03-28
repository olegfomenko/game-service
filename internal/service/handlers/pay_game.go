package handlers

import (
	"encoding/json"
	"github.com/olegfomenko/game-service/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func PayGame(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewPayGame(r)
	if err != nil {
		Log(r).WithError(err).Error("invalid request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	gam, err := Connector(r).Asset(request.Data.Attributes.GameCoinId)
	if err != nil {
		Log(r).WithError(err).Error("error getting gam")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	respTx, err := donate(
		r,
		request.Data.Attributes.OwnerId,
		request.Data.Attributes.GameCoinId,
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
