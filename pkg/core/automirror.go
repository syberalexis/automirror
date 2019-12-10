package core

import (
	"fmt"
	"github.com/syberalexis/automirror/pkg/mirrors"
	"strings"
)

type Automirror struct {
	Config  string
	mirrors []mirrors.Mirror
}

func (a *Automirror) GetMirrors() error {
	var mirrorList []string
	for _, mirror := range a.mirrors {
		mirrorList = append(mirrorList, mirror.Name)
	}
	_, err := fmt.Println("[", strings.Join(mirrorList, ","), "]")
	return err
}

func (a *Automirror) Start(mirror string) error {
	return nil
}

func (a *Automirror) Status(mirror string) error {
	return nil
}

func (a *Automirror) Stop(mirror string) error {
	return nil
}

func (a *Automirror) Restart(mirror string) error {
	return nil
}

//type Automirror struct {
//	config  configs.TomlConfig
//	mirrors []mirrors.Mirror
//}
//
//
//func (a *Automirror) buildMirrors() {
//	var mirrorsArray []mirrors.Mirror
//	for name, mirror := range a.config.Mirrors {
//		var puller pullers.Puller
//		var pusher pushers.Pusher
//
//		if engine := buildEngine(mirror.Puller); engine != nil {
//			puller = engine.(pullers.Puller)
//		}
//		if engine := buildEngine(mirror.Pusher); engine != nil {
//			pusher = engine.(pushers.Pusher)
//		}
//
//		mirrorsArray = append(
//			mirrorsArray,
//			mirrors.New(
//				name,
//				puller,
//				pusher,
//				mirror.Timer,
//				logger.LoggerInfo{
//					Directory: a.config.LogDir,
//					Filename:  name,
//					Format:    a.config.LogFormat,
//					Level:     a.config.LogLevel,
//				},
//			),
//		)
//	}
//	a.mirrors = mirrorsArray
//}
//
//func initializeLogger(config configs.TomlConfig) *os.File {
//	var file *os.File
//	if config.LogDir != "" {
//		file, err := os.OpenFile(filesystem.Combine(config.LogDir, "automirror.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
//		if err != nil {
//			log.Fatal(err)
//		}
//		log.SetOutput(file)
//	}
//	if config.LogFormat == "json" {
//		log.SetFormatter(&log.JSONFormatter{})
//	}
//	if config.LogLevel != "" {
//		var level log.Level
//		ptr := &level
//		err := ptr.UnmarshalText([]byte(config.LogLevel))
//		if err != nil {
//			log.Fatal(err)
//		}
//		log.SetLevel(level)
//	}
//
//	return file
//}
//
//func readConfiguration(configFile string) configs.TomlConfig {
//	var config configs.TomlConfig
//	err := configs.Parse(&config, configFile)
//	if err != nil {
//		log.Fatal(err)
//	}
//	return config
//}
//
//func buildEngine(config configs.EngineConfig) interface{} {
//	var engine interface{}
//	var err error
//
//	switch config.Name {
//	case "deb":
//		engine, err = pullers.NewDeb(config)
//	case "docker":
//		engine, err = both.NewDocker(config)
//	case "git":
//		engine, err = both.NewGit(config)
//	case "mvn":
//		engine, err = pullers.NewMaven(config)
//	case "pip":
//		engine, err = pullers.NewPython(config)
//	case "rsync":
//		engine, err = both.NewRsync(config)
//	case "wget":
//		engine, err = pullers.NewWget(config)
//	case "jfrog":
//		engine, err = pushers.NewJFrog(config)
//	default:
//		engine = nil
//	}
//
//	if err != nil {
//		log.Error(err)
//	}
//
//	return engine
//}
//
//func main() {
//	var configFile = "/etc/automirror/config.toml"
//
//	args := os.Args[1:]
//	if len(args) > 0 {
//		configFile = args[0]
//	}
//
//	// Read configuration
//	automirror := Automirror{config: readConfiguration(configFile)}
//
//	// Logging
//	file := initializeLogger(automirror.config)
//	defer file.Close()
//	defer logger.CloseLoggers()
//
//	// Build mirrors
//	automirror.buildMirrors()
//
//	// Run mirrors
//	var wg sync.WaitGroup
//	for _, mirror := range automirror.mirrors {
//		wg.Add(1)
//		go func(mirror mirrors.Mirror) {
//			defer wg.Done()
//			mirror.Start()
//		}(mirror)
//	}
//	wg.Wait()
//}
