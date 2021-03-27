/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type PayPlayer struct {
	Key
	Attributes PayPlayerAttributes `json:"attributes"`
}
type PayPlayerResponse struct {
	Data     PayPlayer `json:"data"`
	Included Included  `json:"included"`
}

type PayPlayerListResponse struct {
	Data     []PayPlayer `json:"data"`
	Included Included    `json:"included"`
	Links    *Links      `json:"links"`
}

// MustPayPlayer - returns PayPlayer from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustPayPlayer(key Key) *PayPlayer {
	var payPlayer PayPlayer
	if c.tryFindEntry(key, &payPlayer) {
		return &payPlayer
	}
	return nil
}
