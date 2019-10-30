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
	if m.Puller != nil {
		m.Puller.Pull()
	}
	if m.Pusher != nil {
		m.Pusher.Push()
	}
	log.Print(m.Name + " is up to date !")
}
