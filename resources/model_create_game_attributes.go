/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import regources "gitlab.com/tokend/regources/generated"

type CreateGameAttributes struct {
	NameCompetition string                 `json:"name_competition"`
	Price           *regources.Amount      `json:"price,omitempty"`
	Team1           map[string]interface{} `json:"team1"`
	Team2           map[string]interface{} `json:"team2"`
}
