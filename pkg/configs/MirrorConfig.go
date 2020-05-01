package configs

// MirrorConfig structure from configuration file
type MirrorConfig struct {
	Timer  string       `yaml|toml|json:"timer"`
	Puller EngineConfig `yaml|toml|json:"puller"`
	Pusher EngineConfig `yaml|toml|json:"pusher"`
}
