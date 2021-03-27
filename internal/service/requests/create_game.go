package requests

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/olegfomenko/game-service/resources"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
)

func NewCreateGame(r *http.Request) (resources.CreateGameResponse, error) {
	var request resources.CreateGameResponse

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, validateCreateGame(request)
}

//AssetCode       string                 `json:"asset_code"`
//Date            string                 `json:"date"`
//NameCompetition string                 `json:"name_competition"`
//OwnerId         string                 `json:"owner_id"`
//Price           *regources.Amount      `json:"price,omitempty"`
//SourceBalanceId string                 `json:"source_balance_id"`
//Team1           map[string]interface{} `json:"team1"`
//Team2           map[string]interface{} `json:"team2"`
func validateCreateGame(r resources.CreateGameResponse) error {
	return validation.Errors{
		"/data/attributes/new_competition": validation.Validate(r.Data.Attributes.NameCompetition, validation.Required),
		"/data/attributes/price":           validation.Validate(r.Data.Attributes.Price, validation.Required),
		"/data/attributes/team1": validation.Validate(
			r.Data.Attributes.Team1, validation.Required, validation.Length(6, 6)),
		"/data/attributes/team2": validation.Validate(
			r.Data.Attributes.Team2, validation.Required, validation.Length(6, 6)),
		"/data/attributes/asset_code":        validation.Validate(r.Data.Attributes.AssetCode, validation.Required),
		"/data/attributes/date":              validation.Validate(r.Data.Attributes.Date, validation.Required),
		"/data/attributes/owner_id":          validation.Validate(r.Data.Attributes.OwnerId, validation.Required),
		"/data/attributes/source_balance_id": validation.Validate(r.Data.Attributes.SourceBalanceId, validation.Required),
	}.Filter()
}
