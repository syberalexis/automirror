package configs

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
)

type ConfigManager interface {
	Parse(file string) TomlConfig
}

type TomlConfigManager struct {
	File string
}

func (t TomlConfigManager) Parse() TomlConfig {
	var config TomlConfig
	tomlFile, err := ioutil.ReadFile(t.File)

	if err != nil {
		log.Fatal(err)
	}

	if _, err := toml.Decode(string(tomlFile), &config); err != nil {
		log.Fatal(err)
	}

	return config
}
