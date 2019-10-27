package mirrors

import (
	"automirror/pullers"
	"automirror/pushers"
)

type Mirror struct {
	Puller pullers.Puller
	Pusher pushers.Pusher
}

func (m Mirror) Run() {
	m.Puller.Pull()
	m.Pusher.Push()
}
