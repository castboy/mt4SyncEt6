package account_group

import (
	"fmt"
	"github.com/spf13/viper"
)

var grpCfg = viper.New()
var Groups []string

func init() {
	grpCfg.SetConfigName("account_group")
	grpCfg.AddConfigPath(".")
	grpCfg.AddConfigPath("./account_group")

	err := grpCfg.ReadInConfig()
	if err != nil {
		fmt.Errorf("Fatal error config ")
	}

	Groups = grpCfg.GetStringSlice("group")
}
