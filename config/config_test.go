package config

import (
	"os"
	"testing"

	"github.com/pelletier/go-toml/v2"
)

func TestConfig(t *testing.T) {
	var config Config
	data, err := os.ReadFile("Config.toml")
	if err != nil {
		t.Fatal(err)
	}
	err = toml.Unmarshal(data, &config)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v\n", config)
}

func TestGetConfig(t *testing.T) {
	t.Logf("%#v\n", Get())
}
