package config

import (
	"testing"

	"github.com/spf13/viper"
)

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
