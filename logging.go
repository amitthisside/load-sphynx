package main

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func initLogging(vs *VirtualService) *logrus.Logger {
	timestamp := time.Now().Format("20060102_150405")
	logFileName := filepath.Join("logs", vs.Algorithm+"_"+strconv.Itoa(vs.Port)+"_"+timestamp+".log")

	// Ensure the logs directory exists
	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		logrus.Fatalf("Failed to create logs directory: %v", err)
	}

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
