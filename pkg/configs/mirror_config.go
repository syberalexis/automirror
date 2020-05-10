package configs

// MirrorConfig structure from configuration file
type MirrorConfig struct {
	Timer  string       `yaml:"timer"`
	Engine EngineConfig `yaml:"engine"`
}
