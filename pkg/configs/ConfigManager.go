package configs

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
)

// Parse method to read TOML file and convert it into a struct
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
