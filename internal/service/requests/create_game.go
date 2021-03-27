package requests

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/olegfomenko/game-service/resources"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
)

func NewCreateGame (r *http.Request) (resources.CreateGameResponse, error)  {
	var request resources.CreateGameResponse
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, validateCreateGame(request)
}


func validateCreateGame(r resources.CreateGameResponse) error {
	return validation.Errors{
		"/data/attributes/new_competition":       validation.Validate(r.Data.Attributes.NameCompetition, validation.Required),
		"/data/attributes/price":         validation.Validate(r.Data.Attributes.Price, validation.Required),
	}.Filter()
}