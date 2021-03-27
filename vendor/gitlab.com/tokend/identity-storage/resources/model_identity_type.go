/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type IdentityType struct {
	Key
	Attributes    IdentityTypeAttributes    `json:"attributes"`
	Relationships IdentityTypeRelationships `json:"relationships"`
}
type IdentityTypeResponse struct {
	Data     IdentityType `json:"data"`
	Included Included     `json:"included"`
}

type IdentityTypeListResponse struct {
	Data     []IdentityType  `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

func (r *IdentityTypeListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *IdentityTypeListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustIdentityType - returns IdentityType from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustIdentityType(key Key) *IdentityType {
	var identityType IdentityType
	if c.tryFindEntry(key, &identityType) {
		return &identityType
	}
	return nil
}
