package base

import (
	"os"
	"path/filepath"
	"strings"

	"vpn-web.funcworks.net/gb"
)

func initYamlConfig() {
	gb.Viper.AutomaticEnv()
	gb.Viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if path, err := os.Executable(); err != nil {
		gb.Logger.Warn(err.Error())
	} else {
		gb.Viper.AddConfigPath(filepath.Dir(path))
	}

	gb.Viper.AddConfigPath(".")
	gb.Viper.SetConfigName("app")
	gb.Viper.SetConfigType("yaml")

	if err := gb.Viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := gb.Viper.Unmarshal(&gb.Config); err != nil {
		panic(err)
	}

	gb.Logger.Debug("loaded app.config")
}
