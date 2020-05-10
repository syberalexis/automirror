package pullers

import "github.com/syberalexis/automirror/pkg/old"

// Puller interface to expose methods for pulling processes
type Puller interface {
	GetPackages() ([]old.PackageInfo, error)
	GetDependencies() ([]old.PackageInfo, error)
	Pull(name string, version string) error
}
