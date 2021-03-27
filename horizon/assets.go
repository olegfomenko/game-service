package horizon

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	regources "gitlab.com/tokend/regources/generated"
)

func (c *Connector) Asset(code string) (*regources.Asset, error) {
	return c.assetsGetter.get(c, code)
}

type assetsGetter struct{}

func (g *assetsGetter) get(connector *Connector, asset string) (*regources.Asset, error) {
	var response regources.AssetResponse
	err := connector.One("/v3/assets", &assetPathParams{
		AssetCode: asset,
	}).Get(&response)
	if err != nil {
		if err == ErrNotFound {
			return nil, err
		}

		return nil, errors.Wrap(err, "failed to get asset")
	}

	return &response.Data, nil
}

type assetPathParams struct {
	AssetCode string
}

func (p *assetPathParams) Path() string {
	return p.AssetCode
}
