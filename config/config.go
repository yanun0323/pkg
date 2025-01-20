package config

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/yanun0323/pkg/logs"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

/*
Init initial the config from yaml file, and find file from paths relative to where the entry func.

	config.Init("config", true, "../config", "../../config")
*/
func Init(cfgName string, dump bool, relativePaths ...string) error {
	var (
		dir string
		err error
		log logs.Logger
	)

	if dump {
		log = logs.New(logs.LevelInfo)
	}

	sync.OnceFunc(func() {
		dir, err = os.Getwd()
		if err != nil {
			err = errors.Wrap(err, "get wd")
			return
		}
		for _, p := range relativePaths {
			path := filepath.Join(dir, p)
			viper.AddConfigPath(path)
			if dump {
				log.Info("config path: ", path)
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

func InitAndLoad[T any](cfgName string, dump bool, relativePaths ...string) (*T, error) {
	if err := Init(cfgName, dump, relativePaths...); err != nil {
		return nil, err
	}

	var cfg T
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, errors.Wrap(err, "unmarshal config")
	}

	return &cfg, nil
}

func dumpConfig(log logs.Logger) {
	keys := viper.AllKeys()
	sort.Strings(keys)
	for _, key := range keys {
		// if strings.Contains(key, "password") || strings.Contains(key, "secret") || strings.Contains(key, "key") || strings.Contains(key, "pass") || strings.Contains(key, "pem") {
		// 	continue
		// }
		log.Info(key, ": ", viper.Get(key))
	}
}
