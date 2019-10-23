package main

import (
	"automirror/configs"
	"automirror/pullers"
)

func main() {
	configManager := configs.TomlConfigManager{File: "config.toml"}
	config := configManager.Parse()

	var pullerArray []pullers.Puller
	for name, mirror := range config.Mirrors {
		//var puller pullers.Puller
		//puller = reflect.TypeOf(module.Puller).New()
		println(name)
		println(mirror.Name)
		//reflect.ValueOf(&puller).MethodByName(module.Puller).Call([]reflect.Value{})
		//reflect.TypeOf(&puller).MethodByName(module.Puller)
		//puller = new(reflect.TypeOf(&puller))
		pullerArray = append(pullerArray, pullers.Maven{
			Url:              "https://repo1.maven.org/maven2",
			MetadataFileName: "maven-metadata.xml",
			DatabaseFile:     "maven.db",
			Artifacts:        []pullers.Artifact{{"com.airbnb", "deeplinkdispatch", "4.0.0"}},
		})
	}

	for _, puller := range pullerArray {
		puller.Pull()
	}
}
