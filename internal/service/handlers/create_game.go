package handlers

import (
	"encoding/json"
	"github.com/olegfomenko/game-service/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/tokend/go/xdrbuild"
	"net/http"
	"strings"
)

const IssuanceAmount = 1000000

func CreateGame(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewCreateGame(r)
	if err != nil {
		Log(r).WithError(err).Error("invalid request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	// TODO check error
	details, err := json.Marshal(map[string]interface{}{
		"date":  request.Data.Attributes.Date,
		"team1": request.Data.Attributes.Team1,
		"team2": request.Data.Attributes.Team2,
	})

	date := strings.ReplaceAll(request.Data.Attributes.Date, "-", "")
	date = strings.ReplaceAll(date, " ", "")
	date = strings.ReplaceAll(date, ":", "")

	assetCode := "GAM" + date

	assetType, err := Connector(r).GetUint32KeyValue("asset_type:gam")

	// TODO check TrailingDigitsCount
	createGam := &xdrbuild.CreateAsset{
		RequestID:                0,
		Code:                     assetCode,
		MaxIssuanceAmount:        IssuanceAmount,
		PreIssuanceSigner:        Connector(r).Signer().Seed(),
		InitialPreIssuanceAmount: IssuanceAmount,
		TrailingDigitsCount:      6,
		Policies:                 0,
		Type:                     uint64(assetType),
		CreatorDetails:           json.RawMessage(details),
		AllTasks:                 nil,
	}

	// TODO issuer game asset PlayerGam, TeamGam

	// TODO catch error
	gamBalance, err := Connector(r).Balance(Connector(r).Signer().Address(), assetCode)

	// TODo fix reference
	paymentToOrg, _ := xdrbuild.CreatePaymentForAccount(xdrbuild.CreatePaymentForAccountOpts{
		SourceAccountID:      &gamBalance.Relationships.Owner.Data.ID,
		SourceBalanceID:      gamBalance.ID,
		DestinationAccountID: "",
		Amount:               uint64(*request.Data.Attributes.Price),
		Subject:              "Creating game organizer payment",
		Reference:            "creating_game_pay_admin" + assetCode,
		Fee: xdrbuild.Fee{
			SourceFixed:        0,
			SourcePercent:      0,
			DestinationFixed:   0,
			DestinationPercent: 0,
			SourcePaysForDest:  false,
		},
	})

	orgBalance, err := Connector(r).GetBalanceByID(request.Data.Attributes.SourceBalanceId)

	// TODo fix reference
	paymentToAdmin, _ := xdrbuild.CreatePaymentForAccount(xdrbuild.CreatePaymentForAccountOpts{
		SourceAccountID:      &orgBalance.Relationships.Owner.Data.ID,
		SourceBalanceID:      orgBalance.ID,
		DestinationAccountID: gamBalance.Relationships.Owner.Data.ID,
		Amount:               uint64(*request.Data.Attributes.Price),
		Subject:              "Creating game organizer payment",
		Reference:            "creating_game_pay_org" + assetCode,
		Fee: xdrbuild.Fee{
			SourceFixed:        0,
			SourcePercent:      0,
			DestinationFixed:   0,
			DestinationPercent: 0,
			SourcePaysForDest:  false,
		},
	})

	// TODO catch error
	Connector(r).TxSigned(Connector(r).Signer(), createGam, paymentToOrg, paymentToAdmin)
}
