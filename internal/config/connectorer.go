package config

import (
	"github.com/olegfomenko/game-service/internal/horizon"
	"net/http"
	"net/url"

	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/tokend/connectors/signed"
	"gitlab.com/tokend/keypair"
	"gitlab.com/tokend/keypair/figurekeypair"
)

func (c *config) Connector() *horizon.Connector {
	return c.connector.Do(func() interface{} {
		var config struct {
			Endpoint *url.URL        `fig:"endpoint,required"`
			Signer   keypair.Full    `fig:"signer,required"`
			Source   keypair.Address `fig:"source,required"`
		}

		err := figure.
			Out(&config).
			With(figure.BaseHooks, figurekeypair.Hooks).
			From(kv.MustGetStringMap(c.getter, "horizon")).
			Please()
		if err != nil {
			panic(err)
		}

		cli := signed.NewClient(http.DefaultClient, config.Endpoint)
		if config.Signer != nil {
			cli = cli.WithSigner(config.Signer)
		}

		return horizon.NewConnector(cli, config.Source, config.Signer)
	}).(*horizon.Connector)
}
