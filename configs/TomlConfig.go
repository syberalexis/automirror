package configs

// TomlConfig structure from configuration file
type TomlConfig struct {
	LogFile   string                  `toml:"log_file"`
	LogFormat string                  `toml:"log_format"`
	LogLevel  string                  `toml:"log_level"`
	Mirrors   map[string]MirrorConfig `toml:"mirrors"`
}

// MirrorConfig structure from configuration file
type MirrorConfig struct {
	Timer  string
	Puller EngineConfig `toml:"puller"`
	Pusher EngineConfig `toml:"pusher"`
}

// EngineConfig structure from configuration file
type EngineConfig struct {
	Name   string
	Config string
}
