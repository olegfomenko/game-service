/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import regources "gitlab.com/tokend/regources/generated"

type CreateGameAttributes struct {
	Date            string            `json:"date"`
	NameCompetition string            `json:"name_competition"`
	Price           *regources.Amount `json:"price,omitempty"`
	Team1           Team              `json:"team1"`
	Team2           Team              `json:"team2"`
}
