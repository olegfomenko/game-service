/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type PayGame struct {
	Key
	Attributes PayGameAttributes `json:"attributes"`
}
type PayGameResponse struct {
	Data     PayGame  `json:"data"`
	Included Included `json:"included"`
}

type PayGameListResponse struct {
	Data     []PayGame `json:"data"`
	Included Included  `json:"included"`
	Links    *Links    `json:"links"`
}

// MustPayGame - returns PayGame from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustPayGame(key Key) *PayGame {
	var payGame PayGame
	if c.tryFindEntry(key, &payGame) {
		return &payGame
	}
	return nil
}
