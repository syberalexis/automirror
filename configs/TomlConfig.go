package configs

type TomlConfig struct {
	LogFile   string                  `toml:"log_file"`
	LogFormat string                  `toml:"log_format"`
	LogLevel  string                  `toml:"log_level"`
	Mirrors   map[string]MirrorConfig `toml:"mirrors"`
}

type MirrorConfig struct {
	Timer  string
	Puller EngineConfig `toml:"puller"`
	Pusher EngineConfig `toml:"pusher"`
}

type PullerConfig struct {
	Name        string
	Source      string
	Destination string
	Config      string
}

type PusherConfig struct {
	Name   string
	Config string
}

type EngineConfig struct {
	Name   string
	Config string
}
