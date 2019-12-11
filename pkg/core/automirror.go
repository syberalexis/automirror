package core

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/syberalexis/automirror/pkg/both"
	"github.com/syberalexis/automirror/pkg/configs"
	"github.com/syberalexis/automirror/pkg/mirrors"
	"github.com/syberalexis/automirror/pkg/pullers"
	"github.com/syberalexis/automirror/pkg/pushers"
	"github.com/syberalexis/automirror/utils/logs"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Automirror struct {
	ConfigFile string
	config     configs.TomlConfig
	mirrors    map[string]mirrors.Mirror
	running    bool
	waitGroup  sync.WaitGroup
	logFile    os.File
	logger     *log.Logger
}

func (a *Automirror) Init() {
	a.config = readConfiguration(a.ConfigFile)

	// Logging
	a.logFile, a.logger = logs.NewLogger(
		logs.LoggerInfo{
			Directory: a.config.LogDir,
			Filename:  "automirror",
			Format:    a.config.LogFormat,
			Level:     a.config.LogLevel,
		},
	)

	// Build mirrors
	a.buildMirrors()
}

func (a *Automirror) Destroy() {
	a.logFile.Close()
	for _, mirror := range a.mirrors {
		mirror.Destroy()
	}
}

func (a *Automirror) GetMirrors() error {
	var mirrorList []string
	for _, mirror := range a.mirrors {
		mirrorList = append(mirrorList, mirror.Name)
	}
	_, err := fmt.Printf("[ %s ]\n", strings.Join(mirrorList, ", "))
	return err
}

func (a *Automirror) Start() {
	a.running = true

	for name := range a.mirrors {
		if mirror, exist := a.mirrors[name]; exist {
			a.waitGroup.Add(1)
			go func() {
				runtime.LockOSThread() // Trying it
				defer a.waitGroup.Done()
				fmt.Printf("%d\n", os.Getpid())
				mirror.Start()
			}()
		}
	}

	a.waitGroup.Add(1)
	go func() {
		timer, _ := time.ParseDuration("1s")
		defer a.waitGroup.Done()
		for a.running {
			time.Sleep(timer)
		}
	}()

	a.waitGroup.Wait()
}

func (a *Automirror) Status() {
	for name := range a.mirrors {
		err := a.StatusMirror(name)
		if err != nil {
			a.logger.Error(err)
		}
	}
}

func (a *Automirror) Stop() {
	a.running = false
	for name := range a.mirrors {
		err := a.StopMirror(name)
		if err != nil {
			a.logger.Error(err)
		}
	}
}

func (a *Automirror) Restart() {
	a.Stop()
	a.Start()
}

func (a *Automirror) StartMirror(mirrorName string) error {
	if mirror, exist := a.mirrors[mirrorName]; exist {
		a.waitGroup.Add(1)
		go func() {
			defer a.waitGroup.Done()
			mirror.Start()
		}()
		return nil
	}
	a.waitGroup.Wait()
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
		go mirror.Start()
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

		mirrorMap[name] = mirrors.NewMirror(
			name,
			puller,
			pusher,
			mirror.Timer,
			logs.LoggerInfo{
				Directory: a.config.LogDir,
				Filename:  name,
				Format:    a.config.LogFormat,
				Level:     a.config.LogLevel,
			},
		)
	}
	a.mirrors = mirrorMap
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
