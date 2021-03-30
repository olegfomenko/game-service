package handlers

import (
	"github.com/olegfomenko/game-service/internal/service/requests"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/go/xdrbuild"
	"net/http"
	"strings"
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

func parseGameDate(date string) (string, error) {
	arr := strings.Split(date, " ")
	if len(arr) == 0 {
		return "", errors.New("invalid date string")
	}

	arr[0] = strings.ReplaceAll(arr[0], "-", "")
	return arr[0], nil
}
