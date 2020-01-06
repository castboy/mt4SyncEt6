package security

import (
	"fmt"
	"github.com/spf13/viper"
)

type Group struct {
	ID              int    `json:"id" xorm:"id"`
	Name            string `json:"-" xorm:"name"`
	DepositCurrency string `json:"deposit_currency" xorm:"deposit_currency"`
	MarginStopOut   string `json:"margin_stop_out" xorm:"margin_stop_out"`
	HedgeLargeLeg   string `json:"hedge_largeleg" xorm:"margin_mode"`
}

type Security struct {
	ID         int    `json:"-" xorm:"security_id"`
	Enable     string `json:"enable" xorm:"enable_security"`
	Trade      string `json:"trade" xorm:"enable_trade"`
	CommBase   string `json:"comm_base" xorm:"commission"`
	LotMin     string `json:"lot_min" xorm:"lot_min"`
	LotMax     string `json:"lot_max" xorm:"lot_max"`
	LotStep    string `json:"lot_step" xorm:"lot_step"`
	DpreadDiff string `json:"spread_diff" xorm:"spread_diff"`
}

func init() {
	config := viper.New()
	config.SetConfigName("security.json")
	config.AddConfigPath(".")
	config.AddConfigPath("./conf")
	config.WatchConfig()

	err := config.ReadInConfig()
	if err != nil {
		fmt.Println("Fatal error config file:")
	}
}

func GetGroup() error {

	return nil
}
