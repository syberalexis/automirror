package pullers

// Puller interface to expose methods for pulling processes
type Puller interface {
	Pull() (int, error)
}
