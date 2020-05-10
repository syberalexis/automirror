package configs

import "github.com/syberalexis/automirror/pkg/model"

// EngineConfig structure from configuration file
type EngineConfig struct {
	Name     string          `yaml:"name"`
	Config   string          `yaml:"config"`
	Archives []model.Archive `yaml:"archives"`
}
