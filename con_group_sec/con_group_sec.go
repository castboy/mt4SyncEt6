package con_group_sec

import (
	"fmt"
	"github.com/spf13/viper"
)

var secCfg = viper.New()

func init() {
	secCfg.SetConfigName("con_group_sec")
	secCfg.AddConfigPath(".")
	secCfg.AddConfigPath("./con_group_sec")

	err := secCfg.ReadInConfig()
	if err != nil {
		fmt.Errorf("Fatal error config ")
	}
}

func GetSecStringMap(k string) map[string]string {
	v := secCfg.GetStringMapString(k)
	return v
}
