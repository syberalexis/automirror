package main

import (
	"automirror/configs"
	"automirror/mirrors"
	"automirror/pullers"
	"automirror/pushers"
)

func main() {
	configManager := configs.TomlConfigManager{File: "config.toml"}
	config := configManager.Parse()

	var mirrorsArray []mirrors.Mirror
	for name, mirror := range config.Mirrors {
		//var puller pullers.Puller
		//puller = reflect.TypeOf(module.Puller).New()
		println(name)
		println(mirror.Name)
		//reflect.ValueOf(&puller).MethodByName(module.Puller).Call([]reflect.Value{})
		//reflect.TypeOf(&puller).MethodByName(module.Puller)
		//puller = new(reflect.TypeOf(&puller))
		mirrorsArray = append(mirrorsArray, mirrors.Mirror{
			Puller: pullers.Maven{
				Url:              mirror.Puller.Source,
				Folder:           mirror.Puller.Destination,
				MetadataFileName: "maven-metadata.xml",
				DatabaseFile:     "maven.db",
				Artifacts:        []pullers.Artifact{{"com.airbnb", "deeplinkdispatch", "4.0.0"}},
			},
			Pusher: pushers.JFrog{
				Url:    "http://localhost:8081/artifactory/test",
				ApiKey: "AKCp5e2qXnFDWrtw7hJHjjWxR6ei5tCQ3HCvdnSYop6Y8w1vK1GQeUEKeFqSePJXmpCHexcac",
				Source: "tmp/maven",
			},
		})
	}

	for _, mirror := range mirrorsArray {
		mirror.Run()
	}
}
