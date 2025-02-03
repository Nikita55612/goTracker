package config

import (
	"log"
	"os"
	"sync"

	"github.com/pelletier/go-toml/v2"
)

const Path = "Config.toml"

var (
	cfg  *Config
	once sync.Once
)

func Get() *Config {
	once.Do(func() {
		data, err := os.ReadFile(Path)
		if err != nil {
			if err != os.ErrNotExist {
				log.Fatalln("fail to reading config file:", err)
			}
			log.Println("config file is not exist")
			cfg = Default()
			data, _ = toml.Marshal(cfg)
			err = os.WriteFile(Path, data, 0644)
			if err != nil {
				log.Println("fail writing default config to file:  %w\n", err)
			}
			return
		}
		var rcfg Config
		err = toml.Unmarshal(data, &rcfg)
		if err != nil {
			log.Fatalln("parse config file error:", err)
		}
		cfg = &rcfg
	})
	return cfg
}

type Config struct {
	App     App     `toml:"app"`
	Tracker Tracker `toml:"tracker"`
}

type App struct {
	Port            uint16 `toml:"port"`
	DefaultSavePath string `toml:"default_save_path"`
}

type Tracker struct {
	BaseUrl   string `toml:"base_url"`
	ProxyUrl  string `toml:"proxy_url"`
	BBSession string `toml:"bb_session"`
}

func Default() *Config {
	return &Config{
		Tracker: Tracker{
			BaseUrl:   "https://rutracker.org",
			ProxyUrl:  "https://ps1.blockme.site:443",
			BBSession: "0-52335687-cqygg3U3HlXLVNkKPD6R",
		},
	}
}
