package pushers

// Pusher interface to expose methods for pushing processes
type Pusher interface {
	Push() error
}
