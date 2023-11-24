package config

import (
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/yanun0323/pkg/logs"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

/*
Initial the config from config.yaml

	config.Init("config", true, "../config", "../../config")
*/
func Init(cfgName string, dump bool, relativePaths ...string) error {
	var (
		err error
		log logs.Logger
	)

	if dump {
		log = logs.New("config", logs.LevelInfo)
	}

	sync.OnceFunc(func() {
		_, f, _, _ := runtime.Caller(1)
		for _, p := range relativePaths {
			path := filepath.Join(filepath.Dir(f), p)
			viper.AddConfigPath(path)
			if dump {
				log.Info("config path:", path)
			}
		}
		viper.AddConfigPath(".")
		if len(cfgName) == 0 {
			err = errors.New("empty config name")
			return
		}

		viper.SetConfigName(cfgName)
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.SetConfigType("yaml")

		err = viper.ReadInConfig()
		if err != nil {
			err = errors.Wrap(err, "read in config")
			return
		}
		if dump {
			dumpConfig(log)
		}
	})()
	return err
}

func dumpConfig(log logs.Logger) {
	keys := viper.AllKeys()
	sort.Strings(keys)
	for _, key := range keys {
		if strings.Contains(key, "password") || strings.Contains(key, "secret") || strings.Contains(key, "key") || strings.Contains(key, "pass") || strings.Contains(key, "pem") {
			continue
		}
		log.Info(key, ":", viper.Get(key))
	}
}
