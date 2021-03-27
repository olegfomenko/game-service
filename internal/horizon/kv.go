package horizon

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/go/xdr"
	regources "gitlab.com/tokend/regources/generated"
)

func (c *Connector) GetStringKeyValue(key string) (string, error) {
	val, err := c.kv.get(c, key)
	if err != nil {
		return "", errors.Wrap(err, "failed to get kv value")
	}

	if val == nil {
		return "", nil
	}

	if val.Type != xdr.KeyValueEntryTypeString {
		return "", errors.Wrap(err, "value is not of type string")
	}

	if val.Str == nil {
		return "", errors.New("kv value is nil")
	}

	return *val.Str, nil
}

func (c *Connector) GetUint32KeyValue(key string) (uint32, error) {
	val, err := c.kv.get(c, key)
	if err != nil {
		return 0, errors.Wrap(err, "failed to get kv value")
	}

	if val == nil {
		return 0, nil
	}

	if val.Type != xdr.KeyValueEntryTypeUint32 {
		return 0, errors.Wrap(err, "value is not of type string")
	}

	if val.U32 == nil {
		return 0, errors.New("kv value is nil")
	}

	return *val.U32, nil
}

type kvGetter struct{}

type KVPathParams struct {
	Key string
}

func (p *KVPathParams) Path() string {
	return p.Key
}

func (g *kvGetter) get(connector *Connector, key string) (*regources.KeyValueEntryValue, error) {
	var resp regources.KeyValueEntryResponse
	err := connector.One("/v3/key_values", &KVPathParams{
		Key: key,
	}).Get(&resp)
	if err != nil {
		if err == ErrNotFound {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to get key value")
	}

	return &resp.Data.Attributes.Value, nil
}
