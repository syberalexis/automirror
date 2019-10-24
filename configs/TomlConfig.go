package configs

type TomlConfig struct {
	Mirrors map[string]mirror
}

type mirror struct {
	Name   string
	Puller puller
	Pusher pusher
}

type puller struct {
	Name        string
	Source      string
	Destination string
	Config      string
}

type pusher struct {
	Name   string
	Config string
}
