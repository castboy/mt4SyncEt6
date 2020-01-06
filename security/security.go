package security

import (
	"fmt"
	"github.com/spf13/viper"
)

var grpCfg = viper.New()
var secCfg = viper.New()
var Groups []string

func init() {
	grpCfg.SetConfigName("group")
	grpCfg.AddConfigPath(".")
	grpCfg.AddConfigPath("./security")
	secCfg.WatchConfig()

	secCfg.SetConfigName("security")
	secCfg.AddConfigPath(".")
	secCfg.AddConfigPath("./security")
	secCfg.WatchConfig()

	err := grpCfg.ReadInConfig()
	if err != nil {
		fmt.Errorf("Fatal error config ")
	}
	err = secCfg.ReadInConfig()
	if err != nil {
		fmt.Errorf("Fatal error config ")
	}

	Groups = grpCfg.GetStringSlice("group")
}

func GetStringMap(k string) map[string]string {
	v := secCfg.GetStringMapString(k)
	return v
}
