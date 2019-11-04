package pushers

type Pusher interface {
	Push() error
}
