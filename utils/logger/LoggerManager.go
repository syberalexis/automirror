package logger

import (
	log "github.com/sirupsen/logrus"
	"github.com/syberalexis/automirror/utils/filesystem"
	"os"
)

var loggers []os.File

type LoggerInfo struct {
	Directory string
	Filename  string
	Format    string
	Level     string
}

func NewLogger(loggerInfo LoggerInfo) *log.Logger {
	logger := log.New()

	filename := filesystem.Combine(loggerInfo.Directory, loggerInfo.Filename)
	if filename != "" {
		file, err := os.OpenFile(filename+".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			logger.Fatal(err)
		}
		loggers = append(loggers, *file)
		logger.SetOutput(file)
	}

	if loggerInfo.Format == "json" {
		logger.SetFormatter(&log.JSONFormatter{})
	}

	if loggerInfo.Level != "" {
		var level log.Level
		ptr := &level
		err := ptr.UnmarshalText([]byte(loggerInfo.Level))
		if err != nil {
			logger.Fatal(err)
		}
		logger.SetLevel(level)
	}

	return logger
}

func CloseLoggers() {
	for _, file := range loggers {
		file.Close()
	}
}
