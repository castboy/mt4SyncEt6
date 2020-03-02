package security

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var (
	GrpCfg *viper.Viper
	SecCfg *viper.Viper

	Groups        []string
	appConfigPath string
	appMode       string
)

func init() {
	GrpCfg = viper.New()
	SecCfg = viper.New()
	// `APP_MODE` mapping config-path：`dev` -> `dev`, `prod` -> `prod`, `test` -> `test`
	appMode = os.Getenv("APP_MODE")
	if appMode == "" {
		appMode = "dev"
	}
	//File path config
	workPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	appConfigPath = filepath.Join(workPath, "config/resource")

	//file config
	initConfigFromFiles(GrpCfg, "group")
	initConfigFromFiles(SecCfg, "security")

	GrpCfg.SetConfigName("group")
	GrpCfg.AddConfigPath(".")
	GrpCfg.AddConfigPath("./security")

	SecCfg.SetConfigName("security")
	SecCfg.AddConfigPath(".")
	SecCfg.AddConfigPath("./security")

	err = GrpCfg.ReadInConfig()
	if err != nil {
		fmt.Errorf("Fatal error config ")
	}
	err = SecCfg.ReadInConfig()
	if err != nil {
		fmt.Errorf("Fatal error config ")
	}

	Groups = GrpCfg.GetStringSlice("group")
}

func initConfigFromFiles(config *viper.Viper, fileName string) {
	config.SetConfigName(fileName)
	config.AddConfigPath(".")
	config.AddConfigPath(appConfigPath)
	// for modules test
	config.AddConfigPath("./security")
	config.AddConfigPath("../security")
	config.AddConfigPath("../../security")
	config.AddConfigPath("../../../security")
	config.AddConfigPath("../../../../security")
	err := config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file [%s]: %s \n", fileName, err))
	}
}

func GetConfigService(moduleName string) *viper.Viper {
	switch moduleName {
	case "group":
		return GrpCfg.Sub(appMode)
	case "security":
		return SecCfg.Sub(appMode)

	default:
		panic(fmt.Errorf("Unsupported configuration module : %s\n", moduleName))
	}
}
func GetStringMap(k string) map[string]string {
	v := SecCfg.GetStringMapString(k)
	return v
}
