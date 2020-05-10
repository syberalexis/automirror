package configs

type LogConfig struct {
	Dir    string `yaml:"dir"`
	Format string `yaml:"format"`
	Level  string `yaml:"level"`
}
