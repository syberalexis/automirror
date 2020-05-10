package mirrors

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/syberalexis/automirror/pkg/old"
	"github.com/syberalexis/automirror/utils/logs"
	"os"
	"time"
)

// Mirror structure to pull and push packages
type Mirror struct {
	Name       string
	Puller     old.Puller
	Pusher     old.Pusher
	Timer      string
	IsRunning  bool
	LoggerInfo logs.LoggerInfo
	logFile    os.File
	logger     *log.Logger
}

func NewMirror(name string, puller old.Puller, pusher old.Pusher, timer string, loggerInfo logs.LoggerInfo) Mirror {
	return Mirror{
		Name:       name,
		Puller:     puller,
		Pusher:     pusher,
		Timer:      timer,
		IsRunning:  false,
		LoggerInfo: loggerInfo,
	}
}

func (m Mirror) Destroy() error {
	return m.logFile.Close()
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
		//counter := 0
		//packagesInfo, err := m.Puller.GetDependencies()
		//
		//for _, packageInfo := range packagesInfo {
		//
		//}
		//
		//count, err := m.Puller.Pull(m.logger)
		//if err != nil {
		//	m.logger.Errorf("The %s mirror stop to pull (%d elements). This is due to : %s", m.Name, count, err)
		//}
		timer, _ := time.ParseDuration("2m")
		time.Sleep(timer)
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
