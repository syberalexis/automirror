package configs

type LogConfig struct {
	Dir    string `yaml|toml|json:"dir"`
	Format string `yaml|toml|json:"format"`
	Level  string `yaml|toml|json:"level"`
}
