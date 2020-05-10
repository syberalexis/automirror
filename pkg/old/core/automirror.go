package core

import (
	log "github.com/sirupsen/logrus"
	"github.com/syberalexis/automirror/pkg/old"
	"github.com/syberalexis/automirror/utils/logs"
	"os"
	"sync"
)

type Automirror struct {
	ConfigFile string
	config     old.Config
	mirrors    map[string]old.Mirror
	waitGroup  sync.WaitGroup
	logFile    os.File
	logger     *log.Logger
}

func (a *Automirror) Init() {
	a.config = readConfiguration(a.ConfigFile)

	// Logging
	a.logFile, a.logger = logs.NewLogger(
		logs.LoggerInfo{
			Directory: a.config.LogConfig.Dir,
			Filename:  "automirror",
			Format:    a.config.LogConfig.Format,
			Level:     a.config.LogConfig.Level,
		},
	)

	// Build mirrors
	a.buildMirrors()

	a.start()
}

func (a *Automirror) Destroy() {
	err := a.logFile.Close()
	if err != nil {
		log.Error(err)
	}
	for _, mirror := range a.mirrors {
		err = mirror.Destroy()
		if err != nil {
			log.Error(err)
		}
	}
}

func (a *Automirror) start() {
	for name := range a.mirrors {
		if mirror, exist := a.mirrors[name]; exist {
			a.waitGroup.Add(1)
			go func() {
				defer a.waitGroup.Done()
				mirror.Start()
			}()
		}
	}

	a.waitGroup.Wait()
}

func (a *Automirror) buildMirrors() {
	mirrorMap := make(map[string]old.Mirror)
	for name, mirror := range a.config.Mirrors {
		var puller old.Puller
		var pusher old.Pusher

		if engine := buildEngine(mirror.Puller); engine != nil {
			puller = engine.(old.Puller)
		}
		if engine := buildEngine(mirror.Pusher); engine != nil {
			pusher = engine.(old.Pusher)
		}

		mirrorMap[name] = old.NewMirror(
			name,
			puller,
			pusher,
			mirror.Timer,
			logs.LoggerInfo{
				Directory: a.config.LogConfig.Dir,
				Filename:  name,
				Format:    a.config.LogConfig.Format,
				Level:     a.config.LogConfig.Level,
			},
		)
	}
	a.mirrors = mirrorMap
}

func readConfiguration(configFile string) old.Config {
	var config old.Config
	err := old.Parse(&config, configFile)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func buildEngine(config old.EngineConfig) interface{} {
	var engine interface{}
	var err error

	switch config.Name {
	case "deb":
		engine, err = old.NewDeb(config)
	case "docker":
		engine, err = old.NewDocker(config)
	case "git":
		engine, err = old.NewGit(config)
	case "mvn":
		engine, err = old.NewMaven(config)
	case "pip":
		engine, err = old.NewPython(config)
	case "rsync":
		engine, err = old.NewRsync(config)
	case "wget":
		engine, err = old.NewWget(config)
	case "jfrog":
		engine, err = old.NewJFrog(config)
	default:
		engine = nil
	}

	if err != nil {
		log.Error(err)
	}

	return engine
}
