package handlers

import (
	"github.com/olegfomenko/game-service/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/go/xdrbuild"
	"net/http"
)

func CreateGame(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewCreateGame(r)
	if err != nil {
		Log(r).WithError(err).Error("invalid request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	tx := xdrbuild.CreateIssuanceRequest{
		Reference: "",
		Receiver:  "",
		Asset:     "",
		Amount:    0,
		Details:   nil,
		AllTasks:  nil,
	}

	// TODO issuer game asset (Gam), PlayerGam, TeamGam
	// TODO op1: pay Gam admin -> organizer, op2: pay usd org -> admin
}
