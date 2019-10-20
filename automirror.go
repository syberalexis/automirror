package main

import (
	"automirror/pullers"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	var config TomlConfig
	tomlFile, err := ioutil.ReadFile("config.toml")
	check(err)
	if _, err := toml.Decode(string(tomlFile), &config); err != nil {
		log.Fatal(err)
	}

	maven := pullers.Maven{Name: config.Pullers["maven"].Name}
	maven.Pull()
	goo := pullers.Go{}
	goo.Pull()
}

type TomlConfig struct {
	Pullers map[string]puller
}

type puller struct {
	Name string
}
