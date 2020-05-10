package engine

import (
	log "github.com/sirupsen/logrus"
	"github.com/syberalexis/automirror/pkg/configs"
)

type Engine interface {
	Run() error
}

func NewEngine(config configs.EngineConfig, logger *log.Logger) Engine {
	var eng Engine
	var err error

	switch config.Name {
	//case "deb":
	//	engine, err = old.NewDeb(config)
	//case "docker":
	//	engine, err = old.NewDocker(config)
	//case "git":
	//	engine, err = old.NewGit(config)
	//case "mvn":
	//	engine, err = old.NewMaven(config)
	//case "pip":
	//	engine, err = old.NewPython(config)
	//case "rsync":
	//	engine, err = old.NewRsync(config)
	case "wget":
		eng, err = NewWget(config, logger)
	//case "jfrog":
	//	engine, err = old.NewJFrog(config)
	default:
		eng = nil
	}

	if err != nil {
		logger.Error(err)
	}

	return eng
}
