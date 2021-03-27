/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type CreateGame struct {
	Key
	Attributes CreateGameAttributes `json:"attributes"`
}
type CreateGameResponse struct {
	Data     CreateGame `json:"data"`
	Included Included   `json:"included"`
}

type CreateGameListResponse struct {
	Data     []CreateGame `json:"data"`
	Included Included     `json:"included"`
	Links    *Links       `json:"links"`
}

// MustCreateGame - returns CreateGame from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustCreateGame(key Key) *CreateGame {
	var createGame CreateGame
	if c.tryFindEntry(key, &createGame) {
		return &createGame
	}
	return nil
}
