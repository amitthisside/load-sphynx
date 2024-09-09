package main

import (
	"strconv"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func initLogging(vs *VirtualService) *logrus.Logger {
	logFileName := "logs/" + vs.Algorithm + "_" + strconv.Itoa(vs.Port) + ".log"
	logger := logrus.New()
	logger.SetOutput(&lumberjack.Logger{
		Filename:   logFileName,
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	})
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	return logger
}
