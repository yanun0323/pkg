package config

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/spf13/viper"
	"github.com/yanun0323/errors"
)

var once sync.Once

/*
Init initial the config from yaml file, and find file from paths relative to where the entry func.

	config.Init("config", true, "../config", "../../config")
*/
func Init(cfgName string, dump bool, relativePaths ...string) error {
	var (
		dir string
		err error
	)
	if len(cfgName) == 0 {
		return errors.New("empty config name")

	}

	once.Do(func() {
		dir, err = os.Getwd()
		if err != nil {
			err = errors.Errorf("get wd: %+v", err)
			return
		}

		for _, p := range relativePaths {
			path := filepath.Join(dir, p)
			viper.AddConfigPath(path)
			if dump {
				log.Println("config path: ", path)
			}
		}
		viper.AddConfigPath(".")

		viper.SupportedExts = []string{"yaml", "yml"}
		viper.SetConfigName(cfgName)
		viper.SetConfigType("yaml")
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AutomaticEnv()

		err = viper.ReadInConfig()
		if err != nil {
			err = errors.Errorf("read in config: %+v", err)
			return
		}
		if dump {
			dumpConfig()
		}
	})

	return err
}

var (
	store    atomic.Value
	loadLock sync.Mutex
)

// InitAndLoad initializes the config and unmarshals it into a struct.
//
// The config will be cached in memory, so it will be loaded only once.
func InitAndLoad[T any](cfgName string, dump bool, relativePaths ...string) (*T, error) {
	if cfg, ok := store.Load().(*T); ok {
		return cfg, nil
	}

	loadLock.Lock()
	defer loadLock.Unlock()

	if cfg, ok := store.Load().(*T); ok {
		return cfg, nil
	}

	if err := Init(cfgName, dump, relativePaths...); err != nil {
		return nil, err
	}

	var cfg T
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, errors.Errorf("unmarshal config: %+v", err)
	}

	store.Store(&cfg)

	return &cfg, nil
}

func dumpConfig() {
	keys := viper.AllKeys()
	sort.Strings(keys)
	for _, key := range keys {
		// if strings.Contains(key, "password") || strings.Contains(key, "secret") || strings.Contains(key, "key") || strings.Contains(key, "pass") || strings.Contains(key, "pem") {
		// 	continue
		// }
		log.Println(key, ": ", viper.Get(key))
	}
}
