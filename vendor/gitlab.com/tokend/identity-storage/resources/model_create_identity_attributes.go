/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type CreateIdentityAttributes struct {
	AccountId string           `json:"account_id"`
	Details   *json.RawMessage `json:"details,omitempty"`
	Type      string           `json:"type"`
	Value     string           `json:"value"`
}
