package horizon

import (
	"context"
	"github.com/pkg/errors"
	"gitlab.com/tokend/go/xdrbuild"
	regources "gitlab.com/tokend/regources/generated"
)

func (c *Connector) SubmitSigned(ctx context.Context, op ...xdrbuild.Operation) (*regources.TransactionResponse, error) {
	tx, err := c.TxSigned(c.signer, op...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating signed tx")
	}

	base64, err := tx.Marshal()
	if err != nil {
		return nil, errors.Wrap(err, "error marshaling tx")
	}

	respTx, err := c.Submit(ctx, base64, false)
	return respTx, nil
}
