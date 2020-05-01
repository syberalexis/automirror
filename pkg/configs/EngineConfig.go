package configs

// EngineConfig structure from configuration file
type EngineConfig struct {
	Name   string `yaml|toml|json:"name"`
	Config string `yaml|toml|json:"config"`
}
