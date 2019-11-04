package mirrors

import (
	"automirror/pullers"
	"automirror/pushers"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Mirror struct {
	Name      string
	Puller    pullers.Puller
	Pusher    pushers.Pusher
	Timer     string
	IsRunning bool
}

func (m Mirror) Start() {
	timer, _ := time.ParseDuration(m.Timer)
	for {
		var wg sync.WaitGroup
		wg.Add(1)
		go m.Run()
		wg.Wait()
		//runtime.Gosched()
		if m.Timer == "" || timer == 0 {
			break
		}
		time.Sleep(timer)
	}
}

func (m Mirror) Run() {
	if !m.IsRunning {
		log.Infof("%s is running !", m.Name)
		m.IsRunning = true
		if m.Puller != nil {
			count, err := m.Puller.Pull()
			if err != nil {
				log.Errorf("The %s mirror stop to pull (%s elements). This is due to : %s", m.Name, count, err)
			}
		}
		if m.Pusher != nil {
			err := m.Pusher.Push()
			if err != nil {
				log.Errorf("The %s mirror stop to push. This is due to : %s", m.Name, nil)
			}
		}
		m.IsRunning = false
		log.Infof("%s is up to date !", m.Name)
	}
}
