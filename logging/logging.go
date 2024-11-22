package logging

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
)

func SetupLogfile(path string) {
	logWriter := io.MultiWriter(log.Writer(), &lumberjack.Logger{
		Filename: path,
	})
	log.SetOutput(logWriter)
}
