package core

import (
	log "github.com/sirupsen/logrus"
	"github.com/syberalexis/automirror/pkg/configs"
	"github.com/syberalexis/automirror/pkg/engine"
	"github.com/syberalexis/automirror/utils/logs"
	"os"
	"time"
)

type Mirror struct {
	name    string
	timer   time.Duration
	engine  engine.Engine
	logFile os.File
	logger  *log.Logger
}

func NewMirror(name string, config configs.MirrorConfig, loggerInfo logs.LoggerInfo) (Mirror, error) {
	timer, _ := time.ParseDuration(config.Timer)
	logFile, logger := logs.NewLogger(loggerInfo)

	eng := engine.NewEngine(config.Engine, logger)

	return Mirror{
		name:    name,
		timer:   timer,
		engine:  eng,
		logFile: logFile,
		logger:  logger,
	}, nil
}

func (mirror Mirror) Run() {
	defer mirror.logFile.Close()

	for true {
		err := mirror.engine.Run()
		if err != nil {
			mirror.logger.Errorf("Error during mirror's %s running !", mirror.name, err)
		}

		if mirror.timer == 0 {
			break
		}
		time.Sleep(mirror.timer)
	}
}
