package configs

type TomlConfig struct {
	Mirrors   map[string]MirrorConfig `toml:"mirrors"`
	LogFile   string                  `toml:"log_file"`
	LogFormat string                  `toml:"log_format"`
}

type MirrorConfig struct {
	Puller PullerConfig `toml:"puller"`
	Pusher PusherConfig `toml:"pusher"`
	Timer  string
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
