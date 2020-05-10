package configs

// Config structure from configuration file
type Config struct {
	LogConfig LogConfig               `yaml:"log"`
	Mirrors   map[string]MirrorConfig `yaml:"mirrors"`
}
