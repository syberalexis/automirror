package pullers

// Puller interface to expose methods for pulling processes
type Puller interface {
	GetPackages() ([]PackageInfo, error)
	GetDependencies() ([]PackageInfo, error)
	Pull(name string, version string) error
}
