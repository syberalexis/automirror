package main

import (
	"automirror/configs"
	"automirror/mirrors"
	"automirror/pullers"
	"automirror/pushers"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
	"time"
)

func runThread(mirror mirrors.Mirror) {
	log.Print("Mirror " + mirror.Name + " is running")
	for range time.Tick(mirror.Timer * mirror.Unit) {
		mirror.Run()
	}
}

func main() {
	// Read configuration
	var config configs.TomlConfig
	tomlFile, err := ioutil.ReadFile("config.toml")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := toml.Decode(string(tomlFile), &config); err != nil {
		log.Fatal(err)
	}

	// Build mirrors
	var mirrorsArray []mirrors.Mirror
	for name, mirror := range config.Mirrors {
		puller := buildPuller(mirror.Puller)
		pusher := buildPusher(mirror.Pusher)
		unit, _ := time.ParseDuration(string(mirror.Unit))
		mirrorsArray = append(
			mirrorsArray,
			mirrors.Mirror{
				Name:   name,
				Puller: puller,
				Pusher: pusher,
				Timer:  time.Duration(mirror.Timer),
				Unit:   unit,
			},
		)
	}

	// Run mirrors
	for _, mirror := range mirrorsArray {
		go runThread(mirror)
	}
}

func buildPuller(config configs.PullerConfig) pullers.Puller {
	switch config.Name {
	case "maven":
		return pullers.BuildMaven(config)
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
