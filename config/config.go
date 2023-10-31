package config

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/viper"
	"github.com/yanun0323/pkg/logs"
)

/*
Initial the config from config.yaml

	# Sample Code

		var _once sync.Once

		func Init(cfgName string) error {
			var err error
			_once.Do(
				func() {
					_, f, _, _ := runtime.Caller(0)
					cfgPath := filepath.Join(filepath.Dir(f), "../../config")
					if err = config.Init(cfgPath, cfgName, true); err != nil {
						return
					}
				},
			)
			return err
		}
*/
func Init(cfgPath, cfgName string, dump bool) error {
	viper.AddConfigPath(cfgPath)
	viper.AddConfigPath("./config")
	configName := os.Getenv("CONFIG_NAME")
	if len(configName) != 0 {
		cfgName = configName
	}
	if len(cfgName) == 0 {
		cfgName = "config"
	}

	viper.SetConfigName(cfgName)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	if dump {
		Dump()
	}
	return nil
}

func Dump() {
	keys := viper.AllKeys()
	sort.Strings(keys)
	for _, key := range keys {
		if strings.Contains(key, "password") || strings.Contains(key, "secret") || strings.Contains(key, "key") || strings.Contains(key, "pass") || strings.Contains(key, "pem") {
			continue
		}
		logs.New("config", 0).Info(fmt.Sprintf("%s: %+v\n", key, viper.Get(key)))
	}
}
