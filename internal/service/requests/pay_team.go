package requests

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/olegfomenko/game-service/resources"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
)

func NewPayTeam (r *http.Request) (resources.PayTeamResponse, error)  {
	var request resources.PayTeamResponse

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, validatePayTeam(request)
}


func validatePayTeam(r resources.PayTeamResponse) error {
	return validation.Errors{
		"/data/attributes/team_name":       validation.Validate(r.Data.Attributes.TeamName, validation.Required),
		"/data/attributes/amount":      validation.Validate(r.Data.Attributes.Amount, validation.Required),
	}.Filter()
}