package core

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/syberalexis/automirror/pkg/both"
	"github.com/syberalexis/automirror/pkg/configs"
	"github.com/syberalexis/automirror/pkg/mirrors"
	"github.com/syberalexis/automirror/pkg/pullers"
	"github.com/syberalexis/automirror/pkg/pushers"
	"github.com/syberalexis/automirror/utils/filesystem"
	"github.com/syberalexis/automirror/utils/logger"
	"os"
	"strings"
	"sync"
)

type Automirror struct {
	ConfigFile string
	config     configs.TomlConfig
	mirrors    map[string]mirrors.Mirror
	waitGroup  sync.WaitGroup
	logFile    *os.File
}

func (a *Automirror) Init() {
	a.config = readConfiguration(a.ConfigFile)

	// Logging
	a.logFile = initializeLogger(a.config)

	// Build mirrors
	a.buildMirrors()
}

func (a *Automirror) Destroy() {
	a.logFile.Close()
	for _, mirror := range a.mirrors {
		mirror.CloseLogger()
	}
	logger.CloseLoggers()
}

func (a *Automirror) GetMirrors() error {
	var mirrorList []string
	for _, mirror := range a.mirrors {
		mirrorList = append(mirrorList, mirror.Name)
	}
	_, err := fmt.Println("[", strings.Join(mirrorList, ","), "]")
	return err
}

func (a *Automirror) Start() {
	for name := range a.mirrors {
		err := a.StartMirror(name)
		if err != nil {
			log.Fatal(err)
		}
	}
	//waitGroup.Wait()
}

func (a *Automirror) Status() {
	for name := range a.mirrors {
		err := a.StatusMirror(name)
		if err != nil {
			log.Error(err)
		}
	}
}

func (a *Automirror) Stop() {
	a.waitGroup.Done()
}

func (a *Automirror) Restart() {
	a.Stop()
	a.Start()
}

func (a *Automirror) StartMirror(mirrorName string) error {
	if mirror, exist := a.mirrors[mirrorName]; exist {
		a.waitGroup.Add(1)
		go func(mirror mirrors.Mirror) {
			defer a.waitGroup.Done()
			mirror.Start()
		}(mirror)
		a.waitGroup.Wait()
		return nil
	}
	return fmt.Errorf("no mirror found with the name : %s", mirrorName)
}

func (a *Automirror) StatusMirror(mirrorName string) error {
	if mirror, exist := a.mirrors[mirrorName]; exist {
		mirror.Status()
		return nil
	}
	return fmt.Errorf("no mirror found with the name : %s", mirrorName)
}

func (a *Automirror) StopMirror(mirrorName string) error {
	if mirror, exist := a.mirrors[mirrorName]; exist {
		mirror.Stop()
		return nil
	}
	return fmt.Errorf("no mirror found with the name : %s", mirrorName)
}

func (a *Automirror) RestartMirror(mirrorName string) error {
	if mirror, exist := a.mirrors[mirrorName]; exist {
		mirror.Stop()
		a.waitGroup.Add(1)
		go func(mirror mirrors.Mirror) {
			defer a.waitGroup.Done()
			mirror.Start()
		}(mirror)
		a.waitGroup.Wait()
		return nil
	}
	return fmt.Errorf("no mirror found with the name : %s", mirrorName)
}

func (a *Automirror) buildMirrors() {
	mirrorMap := make(map[string]mirrors.Mirror)
	for name, mirror := range a.config.Mirrors {
		var puller pullers.Puller
		var pusher pushers.Pusher

		if engine := buildEngine(mirror.Puller); engine != nil {
			puller = engine.(pullers.Puller)
		}
		if engine := buildEngine(mirror.Pusher); engine != nil {
			pusher = engine.(pushers.Pusher)
		}

		mirrorMap[name] = mirrors.New(
			name,
			puller,
			pusher,
			mirror.Timer,
			logger.LoggerInfo{
				Directory: a.config.LogDir,
				Filename:  name,
				Format:    a.config.LogFormat,
				Level:     a.config.LogLevel,
			},
		)
	}
	a.mirrors = mirrorMap
}

func initializeLogger(config configs.TomlConfig) *os.File {
	var file *os.File
	if config.LogDir != "" {
		file, err := os.OpenFile(filesystem.Combine(config.LogDir, "automirror.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(file)
	}
	if config.LogFormat == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	}
	if config.LogLevel != "" {
		var level log.Level
		ptr := &level
		err := ptr.UnmarshalText([]byte(config.LogLevel))
		if err != nil {
			log.Fatal(err)
		}
		log.SetLevel(level)
	}

	return file
}

func readConfiguration(configFile string) configs.TomlConfig {
	var config configs.TomlConfig
	err := configs.Parse(&config, configFile)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func buildEngine(config configs.EngineConfig) interface{} {
	var engine interface{}
	var err error

	switch config.Name {
	case "deb":
		engine, err = pullers.NewDeb(config)
	case "docker":
		engine, err = both.NewDocker(config)
	case "git":
		engine, err = both.NewGit(config)
	case "mvn":
		engine, err = pullers.NewMaven(config)
	case "pip":
		engine, err = pullers.NewPython(config)
	case "rsync":
		engine, err = both.NewRsync(config)
	case "wget":
		engine, err = pullers.NewWget(config)
	case "jfrog":
		engine, err = pushers.NewJFrog(config)
	default:
		engine = nil
	}

	if err != nil {
		log.Error(err)
	}

	return engine
}
