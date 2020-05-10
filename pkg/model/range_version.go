package model

type RangeVersion struct {
	Min string `yaml:"min"`
	Max string `yaml:"max"`
}

func (rangeVersion RangeVersion) IsInRange(version string) bool {
	isUpToMin := true
	isDownToMax := true

	if len(rangeVersion.Min) > 0 {
		isUpToMin = rangeVersion.Min <= version
	}
	if len(rangeVersion.Max) > 0 {
		isDownToMax = rangeVersion.Max >= version
	}

	return isUpToMin && isDownToMax
}
