package horizon

import (
	"fmt"
	"net/url"
	"strings"

	"gitlab.com/tokend/connectors/signed"

	"gitlab.com/tokend/connectors/lazyinfo"

	"gitlab.com/tokend/connectors/submit"

	horizon "gitlab.com/distributed_lab/json-api-connector"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdrbuild"
	"gitlab.com/tokend/keypair"
)

var (
	ErrDataEmpty = horizon.ErrDataEmpty
	ErrNotFound  = horizon.ErrNotFound
)

type Connector struct {
	*horizon.Connector
	*horizon.Streamer
	*submit.Submitter
	cli *signed.Client

	balances     *balancesGetter
	assetsGetter *assetsGetter
	kv           *kvGetter
	requests     *requestsGetter

	*lazyinfo.LazyInfoer

	source keypair.Address
	signer keypair.Full
}

func NewConnector(cli *signed.Client, source keypair.Address, signer keypair.Full) *Connector {
	return &Connector{
		cli:          cli,
		Submitter:    submit.New(cli),
		Connector:    horizon.NewConnector(cli),
		LazyInfoer:   lazyinfo.New(cli),
		source:       source,
		signer:       signer,
		balances:     &balancesGetter{},
		assetsGetter: &assetsGetter{},
		kv:           &kvGetter{},
		requests:     &requestsGetter{},
	}
}

func (c *Connector) Signer() keypair.Full {
	return c.signer
}

func (c *Connector) Source() keypair.Full {
	return c.signer
}

func (c *Connector) TxHash(envelope string) (string, error) {
	builder, err := c.TXBuilder()
	if err != nil {
		return "", errors.Wrap(err, "failed to init builder")
	}

	return builder.TXHashHex(envelope)
}

func (c *Connector) TxSigned(signer keypair.Full, operations ...xdrbuild.Operation) (*xdrbuild.Transaction, error) {
	tx, err := c.Tx(operations...)
	if err != nil {
		return nil, err
	}
	return tx.Sign(signer), nil
}

func (c *Connector) Tx(operations ...xdrbuild.Operation) (*xdrbuild.Transaction, error) {
	builder, err := c.TXBuilder()
	if err != nil {
		return nil, errors.Wrap(err, "failed to init builder")
	}

	tx := builder.Transaction(c.source)
	for _, op := range operations {
		tx = tx.Op(op)
	}
	return tx, nil
}

func (c *Connector) List(u string, params horizon.Encoder) *horizon.Streamer {
	return horizon.NewStreamer(c.cli, u, params)
}

type pather interface {
	Path() string
}

type Pather struct {
}

func (p *Pather) Path() string {
	return ""
}

type getter struct {
	err error
	c   *Connector
	u   *url.URL
}

func (g *getter) Get(dst interface{}, enc ...horizon.Encoder) error {
	if g.err != nil {
		return g.err
	}

	encodedParams := make([]string, 0, len(enc))
	for _, v := range enc {
		encodedParams = append(encodedParams, v.Encode())
	}
	g.u.RawQuery = strings.Join(encodedParams, "&")

	return g.c.Get(g.u, dst)
}

func (c *Connector) One(basicPath string, path pather) *getter {
	result := getter{
		c: c,
	}
	result.u, result.err = url.Parse(fmt.Sprintf("%s/%s", basicPath, path.Path()))
	return &result
}
