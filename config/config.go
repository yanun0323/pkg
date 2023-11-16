package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/spf13/viper"
	"github.com/yanun0323/pkg/logs"
)

/*
Initial the config from config.yaml

	config.Init("config", true, "../config", "../../config")
*/
func Init(cfgName string, dump bool, relativePaths ...string) error {
	var err error
	sync.OnceFunc(func() {
		_, f, _, _ := runtime.Caller(0)
		for _, p := range relativePaths {
			viper.AddConfigPath(filepath.Join(filepath.Dir(f), p))
		}
		viper.AddConfigPath(".")
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

		err = viper.ReadInConfig()
		if err != nil {
			return
		}
		if dump {
			dumpConfig()
		}
	})
	return err
}

func dumpConfig() {
	keys := viper.AllKeys()
	sort.Strings(keys)
	for _, key := range keys {
		if strings.Contains(key, "password") || strings.Contains(key, "secret") || strings.Contains(key, "key") || strings.Contains(key, "pass") || strings.Contains(key, "pem") {
			continue
		}
		logs.New("config", 0).Info(fmt.Sprintf("%s: %+v\n", key, viper.Get(key)))
	}
}
