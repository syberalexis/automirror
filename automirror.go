package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/syberalexis/automirror/both"
	"github.com/syberalexis/automirror/configs"
	"github.com/syberalexis/automirror/mirrors"
	"github.com/syberalexis/automirror/pullers"
	"github.com/syberalexis/automirror/pushers"
	"os"
	"sync"
)

var configFile = "config.toml"

func main() {
	// Read configuration
	var config configs.TomlConfig
	err := configs.Parse(&config, configFile)
	if err != nil {
		log.Fatal(err)
	}

	// Logging
	if config.LogFile != "" {
		file, err := os.OpenFile(config.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
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
	// Build mirrors
	mirrorsArray := buildMirrors(config)

	// Run mirrors
	var wg sync.WaitGroup
	for _, mirror := range mirrorsArray {
		wg.Add(1)
		go func(mirror mirrors.Mirror) {
			defer wg.Done()
			mirror.Start()
		}(mirror)
	}
	wg.Wait()
}

func buildMirrors(config configs.TomlConfig) []mirrors.Mirror {
	var mirrorsArray []mirrors.Mirror
	for name, mirror := range config.Mirrors {
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
			mirrors.Mirror{
				Name:   name,
				Puller: puller,
				Pusher: pusher,
				Timer:  mirror.Timer,
			},
		)
	}
	return mirrorsArray
}

func buildEngine(config configs.EngineConfig) interface{} {
	var engine interface{}
	var err error

	switch config.Name {
	case "deb":
		engine, err = pullers.NewDeb(config)
	case "docker":
		engine, err = pullers.NewDocker(config)
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
