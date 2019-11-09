package main

import (
	"automirror/configs"
	"automirror/mirrors"
	"automirror/pullers"
	"automirror/pushers"
	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"sync"
)

//var configFile = "/etc/automirror/config.toml"
var configFile = "config.toml"

func main() {
	// Read configuration
	config := readConfiguration()

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
	log.SetLevel(log.InfoLevel)

	// Build mirrors
	mirrorsArray := buildMirrors(config)

	// Run mirrors
	var wg sync.WaitGroup
	for _, mirror := range mirrorsArray {
		wg.Add(1)
		go mirror.Start()
	}
	wg.Wait()
}

func readConfiguration() configs.TomlConfig {
	var config configs.TomlConfig
	tomlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := toml.Decode(string(tomlFile), &config); err != nil {
		log.Fatal(err)
	}
	return config
}

func buildMirrors(config configs.TomlConfig) []mirrors.Mirror {
	var mirrorsArray []mirrors.Mirror
	for name, mirror := range config.Mirrors {
		puller := buildPuller(mirror.Puller)
		pusher := buildPusher(mirror.Pusher)
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

func buildPuller(config configs.PullerConfig) pullers.Puller {
	var puller pullers.Puller
	var err error

	switch config.Name {
	case "apt":
		puller, err = pullers.BuildApt(config)
	case "deb":
		puller, err = pullers.BuildDeb(config)
	case "docker":
		puller, err = pullers.BuildDocker(config)
	case "git":
		puller, err = pullers.BuildGit(config)
	case "mvn":
		puller, err = pullers.BuildMaven(config)
	case "pip":
		puller, err = pullers.BuildPython(config)
	case "rsync":
		puller, err = pullers.BuildRsync(config)
	case "wget":
		puller, err = pullers.BuildWget(config)
	default:
		puller = nil
	}

	if err != nil {
		log.Error(err)
	}

	return puller
}

func buildPusher(config configs.PusherConfig) pushers.Pusher {
	var pusher pushers.Pusher
	var err error

	switch config.Name {
	case "jfrog":
		pusher, err = pushers.BuildJFrog(config)
	//case "rsync":
	//	pusher, err = both.BuildRsync(config)
	default:
		pusher = nil
	}

	if err != nil {
		log.Error(err)
	}

	return pusher
}
