/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type PayTeam struct {
	Key
	Attributes PayTeamAttributes `json:"attributes"`
}
type PayTeamResponse struct {
	Data     PayTeam  `json:"data"`
	Included Included `json:"included"`
}

type PayTeamListResponse struct {
	Data     []PayTeam `json:"data"`
	Included Included  `json:"included"`
	Links    *Links    `json:"links"`
}

// MustPayTeam - returns PayTeam from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustPayTeam(key Key) *PayTeam {
	var payTeam PayTeam
	if c.tryFindEntry(key, &payTeam) {
		return &payTeam
	}
	return nil
}
