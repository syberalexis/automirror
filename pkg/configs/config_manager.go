package configs

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

// Parse method to read config file and convert it into a struct
func Parse(config interface{}, configFile string) error {
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	switch filepath.Ext(configFile) {
	case "json":
		err = json.Unmarshal(file, config)
	case "toml":
		err = toml.Unmarshal(file, config)
	case "yaml":
	case "yml":
	default:
		err = yaml.Unmarshal(file, config)
	}

	if err != nil {
		return err
	}

	return nil
}

func ReadFile(configFile string) (Config, error) {
	var config Config
	err := Parse(&config, configFile)
	return config, err
}
