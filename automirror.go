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
	file, err := os.OpenFile(config.LogFile, os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	//log.SetOutput(file)
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
	switch config.Name {
	case "mvn":
		return pullers.BuildMaven(config)
	case "pip":
		return pullers.BuildPython(config)
	//case "apt":
	//	return pullers.BuildApt(config)
	default:
		return nil
	}
}

func buildPusher(config configs.PusherConfig) pushers.Pusher {
	switch config.Name {
	case "jfrog":
		return pushers.BuildJFrog(config)
	default:
		return nil
	}
}
