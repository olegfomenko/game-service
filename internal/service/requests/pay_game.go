package requests

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/olegfomenko/game-service/resources"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
)

func NewPayGame (r *http.Request) (resources.PayGameResponse, error)  {
	var request resources.PayGameResponse

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, validatePayGame(request)
}


func validatePayGame(r resources.PayGameResponse) error {
	return validation.Errors{
		"/data/attributes/game_coin_id":       validation.Validate(r.Data.Attributes.GameCoinId, validation.Required),
		"/data/attributes/amount":      validation.Validate(r.Data.Attributes.Amount, validation.Required),
	}.Filter()
}