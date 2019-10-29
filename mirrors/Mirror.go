package mirrors

import (
	"automirror/pullers"
	"automirror/pushers"
	"log"
	"time"
)

type Mirror struct {
	Name   string
	Puller pullers.Puller
	Pusher pushers.Pusher
	Timer  time.Duration
	Unit   time.Duration
}

func (m Mirror) Run() {
	m.Puller.Pull()
	m.Pusher.Push()
	log.Print("Mirror " + m.Name + " was up to date !")
}
