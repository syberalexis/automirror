package mirrors

import (
	log "github.com/sirupsen/logrus"
	"github.com/syberalexis/automirror/pullers"
	"github.com/syberalexis/automirror/pushers"
	"github.com/syberalexis/automirror/utils"
	"time"
)

// Mirror structure to pull and push packages
type Mirror struct {
	Name      string
	Puller    pullers.Puller
	Pusher    pushers.Pusher
	Timer     string
	IsRunning bool
	logger    *log.Logger
}

func New(name string, puller pullers.Puller, pusher pushers.Pusher, timer string, loggerInfo utils.LoggerInfo) Mirror {
	return Mirror{
		Name:      name,
		Puller:    puller,
		Pusher:    pusher,
		Timer:     timer,
		IsRunning: false,
		logger:    utils.NewLogger(loggerInfo),
	}
}

// Start method to initialize scheduler
func (m Mirror) Start() {
	timer, _ := time.ParseDuration(m.Timer)
	for {
		m.Run()
		if m.Timer == "" || timer == 0 {
			break
		}
		time.Sleep(timer)
	}
	m.logger.Infof("%s is stopped !", m.Name)
	log.Infof("%s is stopped !", m.Name)
}

// Run method to pull and push if not already running
func (m Mirror) Run() {
	if !m.IsRunning {
		m.logger.Infof("%s is running !", m.Name)
		m.IsRunning = true
		if m.Puller != nil {
			count, err := m.Puller.Pull(m.logger)
			if err != nil {
				m.logger.Errorf("The %s mirror stop to pull (%d elements). This is due to : %s", m.Name, count, err)
				log.Errorf("The %s mirror stop to pull (%d elements). This is due to : %s", m.Name, count, err)
			}
		}
		if m.Pusher != nil {
			err := m.Pusher.Push()
			if err != nil {
				m.logger.Errorf("The %s mirror stop to push. This is due to : %s", m.Name, err)
				log.Errorf("The %s mirror stop to push. This is due to : %s", m.Name, err)
			}
		}
		m.IsRunning = false
		m.logger.Infof("%s is up to date !", m.Name)
		log.Infof("%s is up to date !", m.Name)
	}
}
