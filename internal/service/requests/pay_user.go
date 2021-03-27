package requests

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/olegfomenko/game-service/resources"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
)

func NewPayUser (r *http.Request) (resources.PayUserResponse, error)  {
	var request resources.PayUserResponse

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, validatePayUser(request)
}


func validatePayUser(r resources.PayUserResponse) error {
	return validation.Errors{
		"/data/attributes/user_acc_id":       validation.Validate(r.Data.Attributes.UserAccId, validation.Required),
		"/data/attributes/amount":      validation.Validate(r.Data.Attributes.Amount, validation.Required),
	}.Filter()
}