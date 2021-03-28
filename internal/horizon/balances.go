package horizon

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/url"

	regources "gitlab.com/tokend/regources/generated"
)

func (c *Connector) Balance(account, asset string) (*regources.Balance, error) {
	response, err := c.balances.get(c, account, asset)

	if err != nil {
		return nil, err
	}

	if len(response.Data) == 0 {
		return nil, ErrNotFound
	}

	return &response.Data[0], err
}

func (c *Connector) BalanceWithState(account, asset string) (*regources.Balance, *regources.BalanceState, error) {
	response, err := c.balances.get(c, account, asset)
	if err != nil {
		return nil, nil, err
	}

	if len(response.Data) == 0 {
		return nil, nil, ErrDataEmpty
	}

	balance := response.Data[0]
	state := response.Included.MustBalanceState(balance.Relationships.State.Data.GetKey())

	return &balance, state, err
}

func (c *Connector) GetBalanceByID(balanceID string) (*regources.Balance, error) {
	return c.balances.getByID(c, balanceID)
}

type balancesGetter struct{}

func (g *balancesGetter) get(connector *Connector, account, asset string) (*regources.BalanceListResponse, error) {
	var response regources.BalanceListResponse
	err := connector.List("/v3/balances", &balancesFilters{
		Asset: &asset,
		Owner: &account,

		IncludeState: true,
	}).Next(&response)
	if err != nil {
		if err == ErrNotFound || err == ErrDataEmpty {
			return nil, err
		}
		return nil, errors.Wrap(err, "failed to get balance")
	}

	if len(response.Data) > 1 {
		return nil, errors.New("got several balances for one asset")
	}

	return &response, nil
}

type balancesFilters struct {
	Asset *string
	Owner *string

	IncludeState bool
}

func (f balancesFilters) Encode() string {
	u := url.Values{}

	if f.Asset != nil {
		u.Add("filter[asset]", *f.Asset)
	}
	if f.Owner != nil {
		u.Add("filter[owner]", *f.Owner)
	}

	if f.IncludeState {
		u.Add("include", "state")
	}

	return u.Encode()
}

func (g *balancesGetter) getByID(connector *Connector, balanceID string) (*regources.Balance, error) {
	var response regources.BalanceResponse
	err := connector.One("/v3/balances", &balancePathParams{
		BalanceID: balanceID,
	}).Get(&response)
	if err != nil {
		if err == ErrNotFound {
			return nil, nil
		}

		return nil, err
	}

	return &response.Data, nil
}

type balancePathParams struct {
	BalanceID string
}

func (p *balancePathParams) Path() string {
	return p.BalanceID
}
