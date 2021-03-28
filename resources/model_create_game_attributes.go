/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import regources "gitlab.com/tokend/regources/generated"

type CreateGameAttributes struct {
	Amount          regources.Amount       `json:"amount"`
	AssetCode       string                 `json:"asset_code"`
	Date            string                 `json:"date"`
	NameCompetition string                 `json:"name_competition"`
	OwnerId         string                 `json:"owner_id"`
	PaymentTx       string                 `json:"payment_tx"`
	SourceBalanceId string                 `json:"source_balance_id"`
	StreamLink      string                 `json:"stream_link"`
	Team1           map[string]interface{} `json:"team1"`
	Team2           map[string]interface{} `json:"team2"`
}
