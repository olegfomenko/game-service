/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type PayUser struct {
	Key
	Attributes PayUserAttributes `json:"attributes"`
}
type PayUserResponse struct {
	Data     PayUser  `json:"data"`
	Included Included `json:"included"`
}

type PayUserListResponse struct {
	Data     []PayUser `json:"data"`
	Included Included  `json:"included"`
	Links    *Links    `json:"links"`
}

// MustPayUser - returns PayUser from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustPayUser(key Key) *PayUser {
	var payUser PayUser
	if c.tryFindEntry(key, &payUser) {
		return &payUser
	}
	return nil
}
