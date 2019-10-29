package configs

type TomlConfig struct {
	Mirrors map[string]MirrorConfig `toml:"mirrors"`
}

type MirrorConfig struct {
	Puller PullerConfig `toml:"puller"`
	Pusher PusherConfig `toml:"pusher"`
	Timer  int
	Unit   string
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
