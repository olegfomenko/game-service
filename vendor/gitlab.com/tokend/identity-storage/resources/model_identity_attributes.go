/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type IdentityAttributes struct {
	Details *json.RawMessage `json:"details,omitempty"`
	Hash    string           `json:"hash"`
	Salt    string           `json:"salt"`
	Status  string           `json:"status"`
	Value   string           `json:"value"`
}
