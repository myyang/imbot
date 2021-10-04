package config

// proxy github.com/spf13/viper as global config
// read 'im.yaml' by default, or change configure file to what you want.

import (
	"fmt"

	"github.com/spf13/viper"
)

var Config *viper.Viper

const defaultCfgName = "im.yaml"

func init() {
	g := viper.New()
	g.SetConfigFile(defaultCfgName)
	err := g.ReadInConfig()
	if err != nil {
		fmt.Printf(
			"WARNING:\n  unable to read default config: '%v'\n  use config.Set for valid runtime config\n",
			defaultCfgName,
		)
	} else {
		Config = g
	}
}

func Set(cfg *viper.Viper) { Config = cfg }
