package configs

type TomlConfig struct {
	Timer   string
	Mirrors map[string]MirrorConfig `toml:"mirrors"`
}

type MirrorConfig struct {
	Puller PullerConfig `toml:"puller"`
	Pusher PusherConfig `toml:"pusher"`
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
