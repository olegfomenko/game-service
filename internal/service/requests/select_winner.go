package requests

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/olegfomenko/game-service/resources"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
)

func NewSelectWinner (r *http.Request) (resources.SelectWinnerResponse, error)  {
	var request resources.SelectWinnerResponse

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, validateSelectWinner(request)
}


func validateSelectWinner(r resources.SelectWinnerResponse) error {
	return validation.Errors{
		"/data/attributes/game_coin_id":       validation.Validate(r.Data.Attributes.GameCoinId, validation.Required),
		"/data/attributes/team_name":      validation.Validate(r.Data.Attributes.TeamName, validation.Required),
	}.Filter()
}