package mirrors

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/syberalexis/automirror/pkg/pullers"
	"github.com/syberalexis/automirror/pkg/pushers"
	"github.com/syberalexis/automirror/utils/logs"
	"os"
	"time"
)

// Mirror structure to pull and push packages
type Mirror struct {
	Name       string
	Puller     pullers.Puller
	Pusher     pushers.Pusher
	Timer      string
	IsRunning  bool
	LoggerInfo logs.LoggerInfo
	logFile    os.File
	logger     *log.Logger
}

func NewMirror(name string, puller pullers.Puller, pusher pushers.Pusher, timer string, loggerInfo logs.LoggerInfo) Mirror {
	return Mirror{
		Name:       name,
		Puller:     puller,
		Pusher:     pusher,
		Timer:      timer,
		IsRunning:  false,
		LoggerInfo: loggerInfo,
	}
}

func (m Mirror) Destroy() {
	m.logFile.Close()
}

// Start method to initialize scheduler
func (m Mirror) Start() {
	m.logFile, m.logger = logs.NewLogger(m.LoggerInfo)
	timer, _ := time.ParseDuration(m.Timer)
	m.IsRunning = true
	for m.IsRunning {
		m.run()
		if m.Timer == "" || timer == 0 {
			break
		}
		time.Sleep(timer)
	}
	m.logger.Infof("%s is stopped !", m.Name)
}

// Start method to initialize scheduler
func (m Mirror) Status() {
	runningText := "Stopped"
	if m.IsRunning {
		runningText = "Running"
	}
	fmt.Printf("Mirror %s is %s", m.Name, runningText)
}

// Start method to initialize scheduler
func (m Mirror) Stop() {
	m.IsRunning = false
}

// Start method to initialize scheduler
func (m Mirror) Restart() {

}

// run method to pull and push if not already running
func (m Mirror) run() {
	m.logger.Infof("%s is running !", m.Name)
	m.IsRunning = true
	if m.Puller != nil {
		count, err := m.Puller.Pull(m.logger)
		if err != nil {
			m.logger.Errorf("The %s mirror stop to pull (%d elements). This is due to : %s", m.Name, count, err)
		}
	}
	if m.Pusher != nil {
		err := m.Pusher.Push()
		if err != nil {
			m.logger.Errorf("The %s mirror stop to push. This is due to : %s", m.Name, err)
		}
	}
	m.IsRunning = false
	m.logger.Infof("%s is up to date !", m.Name)
}
