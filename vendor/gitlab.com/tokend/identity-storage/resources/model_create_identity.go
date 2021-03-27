/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type CreateIdentity struct {
	Key
	Attributes CreateIdentityAttributes `json:"attributes"`
}
type CreateIdentityResponse struct {
	Data     CreateIdentity `json:"data"`
	Included Included       `json:"included"`
}

type CreateIdentityListResponse struct {
	Data     []CreateIdentity `json:"data"`
	Included Included         `json:"included"`
	Links    *Links           `json:"links"`
	Meta     json.RawMessage  `json:"meta,omitempty"`
}

func (r *CreateIdentityListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *CreateIdentityListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustCreateIdentity - returns CreateIdentity from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustCreateIdentity(key Key) *CreateIdentity {
	var createIdentity CreateIdentity
	if c.tryFindEntry(key, &createIdentity) {
		return &createIdentity
	}
	return nil
}
