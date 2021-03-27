/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type SelectWinner struct {
	Key
	Attributes SelectWinnerAttributes `json:"attributes"`
}
type SelectWinnerResponse struct {
	Data     SelectWinner `json:"data"`
	Included Included     `json:"included"`
}

type SelectWinnerListResponse struct {
	Data     []SelectWinner `json:"data"`
	Included Included       `json:"included"`
	Links    *Links         `json:"links"`
}

// MustSelectWinner - returns SelectWinner from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustSelectWinner(key Key) *SelectWinner {
	var selectWinner SelectWinner
	if c.tryFindEntry(key, &selectWinner) {
		return &selectWinner
	}
	return nil
}
