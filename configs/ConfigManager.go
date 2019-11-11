package configs

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
)

func Parse(config interface{}, configFile string) error {
	tomlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	if _, err := toml.Decode(string(tomlFile), config); err != nil {
		return err
	}

	return nil
}
