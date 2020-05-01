package configs

// Config structure from configuration file
type Config struct {
	LogConfig LogConfig               `yaml|toml|json:"log"`
	Mirrors   map[string]MirrorConfig `yaml|toml|json:"mirrors"`
}
