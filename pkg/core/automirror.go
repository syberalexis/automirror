package core

import (
	log "github.com/sirupsen/logrus"
	"github.com/syberalexis/automirror/pkg/configs"
	"github.com/syberalexis/automirror/utils/logs"
	"os"
	"sync"
)

type Automirror struct {
	config    configs.Config
	mirrors   map[string]Mirror
	waitGroup sync.WaitGroup
	logFile   os.File
	logger    *log.Logger
}

func NewAutomirror(configFile string) *Automirror {
	config, err := configs.ReadFile(configFile)

	if err != nil {
		log.Fatal(err)
	}

	logFile, logger := logs.NewLogger(
		logs.LoggerInfo{
			Directory: config.LogConfig.Dir,
			Filename:  "automirror",
			Format:    config.LogConfig.Format,
			Level:     config.LogConfig.Level,
		},
	)

	var mirrors = make(map[string]Mirror)
	for mirrorName, mirrorConfig := range config.Mirrors {
		mirror, err := NewMirror(
			mirrorName,
			mirrorConfig,
			logs.LoggerInfo{
				Directory: config.LogConfig.Dir,
				Filename:  mirrorName,
				Format:    config.LogConfig.Format,
				Level:     config.LogConfig.Level,
			},
		)

		if err != nil {
			log.Errorf("Unable to create mirror %s", mirrorName, err)
		}

		mirrors[mirrorName] = mirror
	}

	return &Automirror{
		config:  config,
		logFile: logFile,
		logger:  logger,
		mirrors: mirrors,
	}
}

func (automirror *Automirror) Run() {
	for name := range automirror.mirrors {
		if mirror, exist := automirror.mirrors[name]; exist {
			automirror.waitGroup.Add(1)
			go func() {
				defer automirror.waitGroup.Done()
				mirror.Run()
			}()
		}
	}

	automirror.waitGroup.Wait()
}
