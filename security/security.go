package security

import "github.com/thedevsaddam/gojsonq"

type Group struct {
	ID              int    `json:"id" xorm:"id"`
	Name  			string `json:"-" xorm:"name"`
	DepositCurrency string `json:"deposit_currency" xorm:"deposit_currency"`
	MarginStopOut   string `json:"margin_stop_out" xorm:"margin_stop_out"`
	HedgeLargeLeg   string `json:"hedge_largeleg" xorm:"margin_mode"`
}

func getGroup(msg string) Group{
	groupInfo := Group{}

}