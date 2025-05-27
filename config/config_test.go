package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/yanun0323/pkg/test"
)

type TestConfig struct {
	Test      string `yaml:"test"`
	TestSnake string `yaml:"testSnake"`
	TestCamel string `yaml:"testCamel"`
}

func TestInit(t *testing.T) {
	err := Init("config_test", true, "../config")
	test.RequireNoError(t, err)
	test.RequireEqual(t, "hello", viper.GetString("test"))
}

func TestInitAndLoad(t *testing.T) {
	conf, err := InitAndLoad[TestConfig]("config_test", true)
	test.RequireNoError(t, err)
	test.RequireEqual(t, "hello", viper.GetString("test"))
	test.RequireEqual(t, "hello", conf.Test)
	test.RequireEqual(t, "", conf.TestSnake)
	test.RequireEqual(t, "camel should be success", conf.TestCamel)

	conf, ok := store.Load().(*TestConfig)
	test.RequireTrue(t, ok)
	test.RequireEqual(t, "hello", viper.GetString("test"))
	test.RequireEqual(t, "hello", conf.Test)
	test.RequireEqual(t, "", conf.TestSnake)
	test.RequireEqual(t, "camel should be success", conf.TestCamel)

	conf, err = InitAndLoad[TestConfig]("config_test", true)
	test.RequireNoError(t, err)
	test.RequireEqual(t, "hello", viper.GetString("test"))
	test.RequireEqual(t, "hello", conf.Test)
	test.RequireEqual(t, "", conf.TestSnake)
	test.RequireEqual(t, "camel should be success", conf.TestCamel)
}
