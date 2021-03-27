package requests

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/olegfomenko/game-service/resources"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
)

func NewPayPlayer(r *http.Request) (resources.PayPlayerResponse, error)  {
	var request resources.PayPlayerResponse

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, validatePayPlayer(request)
}


func validatePayPlayer(r resources.PayPlayerResponse) error {
	return validation.Errors{
		"/data/attributes/owner_id":       validation.Validate(r.Data.Attributes.OwnerId, validation.Required),
		"/data/attributes/amount":      validation.Validate(r.Data.Attributes.Amount, validation.Required),
	}.Filter()
}