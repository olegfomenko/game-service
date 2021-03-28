package horizon

import (
	"github.com/pkg/errors"
	regources "gitlab.com/tokend/regources/generated"
	"net/url"
	"strconv"
)

type requestsGetter struct{}

func (c *Connector) GetRedemptionRequests(requestType int, state int) ([]regources.ReviewableRequest, error) {
	return c.requests.get(c, requestType, state)
}

func (g *requestsGetter) get(connector *Connector, requestType int, state int) ([]regources.ReviewableRequest, error) {
	var response regources.ReviewableRequestListResponse

	err := connector.List("/v3/requests", &requestsFilters{
		State: &state,
		Type:  &requestType,
	}).Next(&response)

	if err != nil {
		if err == ErrNotFound || err == ErrDataEmpty {
			return nil, err
		}
		return nil, errors.Wrap(err, "failed to get balance")
	}

	return response.Data, nil
}

type requestsFilters struct {
	State *int
	Type  *int
}

func (f requestsFilters) Encode() string {
	u := url.Values{}

	if f.State != nil {
		u.Add("filter[state]", strconv.Itoa(*f.State))
	}

	if f.Type != nil {
		u.Add("filter[type]", strconv.Itoa(*f.Type))
	}

	return u.Encode()
}
