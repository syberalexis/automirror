package main

import (
	"automirror/configs"
	"automirror/mirrors"
	"automirror/pullers"
	"automirror/pushers"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
	"runtime"
	"time"
)

//var configFile = "/etc/automirror/config.toml"
var configFile = "config.toml"

func main() {
	// Read configuration
	var config configs.TomlConfig
	tomlFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := toml.Decode(string(tomlFile), &config); err != nil {
		log.Fatal(err)
	}
	timer, _ := time.ParseDuration(config.Timer)

	// Build mirrors
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
			},
		)
	}

	// Run mirrors
	for {
		for _, mirror := range mirrorsArray {
			log.Print(mirror.Name + " is running")
			go mirror.Run()
		}
		runtime.Gosched()
		time.Sleep(timer)
	}
}

func buildPuller(config configs.PullerConfig) pullers.Puller {
	switch config.Name {
	case "mvn":
		return pullers.BuildMaven(config)
	case "pip":
		return pullers.BuildPython(config)
	case "apt":
		return pullers.BuildApt(config)
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
