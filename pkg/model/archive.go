package model

type Archive struct {
	Name    string       `yaml:"name"`
	Version string       `yaml:"version"`
	Range   RangeVersion `yaml:"range"`
}
