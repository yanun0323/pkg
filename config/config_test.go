package config

import (
	"testing"

	"github.com/spf13/viper"
)

type TestConfig struct {
	Test      string `yaml:"test"`
	TestSnake string `yaml:"testSnake"`
	TestCamel string `yaml:"testCamel"`
}

func TestInit(t *testing.T) {
	err := Init("config_test", true)
	if err != nil {
		t.Fatal(err)
	}

	got := viper.GetString("test")
	expected := "hello"
	if got != expected {
		t.Fatalf("config mismatch. expected %s, but got %s", expected, got)
	}
}

func TestInitAndLoad(t *testing.T) {
	conf, err := InitAndLoad[TestConfig]("config_test", true)
	if err != nil {
		t.Fatal(err)
	}

	got := viper.GetString("test")
	expected := "hello"
	if got != expected {
		t.Fatalf("config mismatch. expected %s, but got %s", expected, got)
	}

	if conf.Test != "hello" {
		t.Fatalf("config mismatch. expected %s, but got %s", "hello", conf.Test)
	}

	if conf.TestSnake != "" {
		t.Fatalf("config mismatch. expected empty, but got %s", conf.TestSnake)
	}

	if conf.TestCamel != "camel" {
		t.Fatalf("config mismatch. expected %s, but got %s", "camel", conf.TestCamel)
	}
}
