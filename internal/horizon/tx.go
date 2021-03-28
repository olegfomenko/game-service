package horizon

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/tokend/go/xdrbuild"
	"gitlab.com/tokend/keypair"
	regources "gitlab.com/tokend/regources/generated"
)

func (c *Connector) SubmitSigned(ctx context.Context, source keypair.Address, op ...xdrbuild.Operation) (*regources.TransactionResponse, error) {
	builder, err := c.TXBuilder()
	if err != nil {
		return nil, errors.Wrap(err, "error creating  builder")
	}

	if source == nil {
		source = c.source
	}

	tx := builder.Transaction(source)

	for _, operation := range op {
		tx.Op(operation)
	}

	tx.Sign(c.signer)

	base64, err := tx.Marshal()
	if err != nil {
		return nil, errors.Wrap(err, "error marshaling tx")
	}

	fmt.Println(base64)

	respTx, err := c.Submit(ctx, base64, true)
	if err != nil {
		return nil, errors.Wrap(err, "error executing tx")
	}

	return respTx, nil
}
