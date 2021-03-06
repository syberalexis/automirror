package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/syberalexis/automirror/both"
	"github.com/syberalexis/automirror/configs"
	"github.com/syberalexis/automirror/mirrors"
	"github.com/syberalexis/automirror/pullers"
	"github.com/syberalexis/automirror/pushers"
	"github.com/syberalexis/automirror/utils"
	"os"
	"sync"
)

type automirror struct {
	config  configs.TomlConfig
	mirrors []mirrors.Mirror
}

func (a *automirror) buildMirrors() {
	var mirrorsArray []mirrors.Mirror
	for name, mirror := range a.config.Mirrors {
		var puller pullers.Puller
		var pusher pushers.Pusher

		if engine := buildEngine(mirror.Puller); engine != nil {
			puller = engine.(pullers.Puller)
		}
		if engine := buildEngine(mirror.Pusher); engine != nil {
			pusher = engine.(pushers.Pusher)
		}

		mirrorsArray = append(
			mirrorsArray,
			mirrors.New(
				name,
				puller,
				pusher,
				mirror.Timer,
				utils.LoggerInfo{
					Directory: a.config.LogDir,
					Filename:  name,
					Format:    a.config.LogFormat,
					Level:     a.config.LogLevel,
				},
			),
		)
	}
	a.mirrors = mirrorsArray
}

func initializeLogger(config configs.TomlConfig) *os.File {
	var file *os.File
	if config.LogDir != "" {
		file, err := os.OpenFile(utils.Combine(config.LogDir, "automirror.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
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

func main() {
	var configFile = "/etc/automirror/config.toml"

	args := os.Args[1:]
	if len(args) > 0 {
		configFile = args[0]
	}

	// Read configuration
	automirror := automirror{config: readConfiguration(configFile)}

	// Logging
	file := initializeLogger(automirror.config)
	defer file.Close()
	defer utils.CloseLoggers()

	// Build mirrors
	automirror.buildMirrors()

	// Run mirrors
	var wg sync.WaitGroup
	for _, mirror := range automirror.mirrors {
		wg.Add(1)
		go func(mirror mirrors.Mirror) {
			defer wg.Done()
			mirror.Start()
		}(mirror)
	}
	wg.Wait()
}
