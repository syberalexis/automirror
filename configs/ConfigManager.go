package configs

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
)

type ConfigManager interface {
	Parse(config interface{})
}

type TomlConfigManager struct {
	File string
}

func (t TomlConfigManager) Parse(config interface{}) {
	tomlFile, err := ioutil.ReadFile(t.File)

	if err != nil {
		log.Fatal(err)
	}
	if _, err := toml.Decode(string(tomlFile), &config); err != nil {
		log.Fatal(err)
	}
}
