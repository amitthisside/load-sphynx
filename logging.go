package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func initLogging() {
	// Create log directory if it doesn't exist
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err := os.Mkdir("logs", 0755)
		if err != nil {
			log.Fatalf("Failed to create log directory: %v", err)
		}
	}

	// Create a new log file for the session
	logFileName := fmt.Sprintf("logs/session_%s.log", time.Now().Format("20060102_150405"))
	var err error
	logFile, err = os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Set log output to the log file
	log.SetOutput(logFile)
}

// Helper function to log the health changes of servers
func logServerHealthChanges() {
	for _, server := range serverList {
		previousHealth, exists := healthState[server.URL]
		if !exists || previousHealth != server.Health {
			if server.Health {
				log.Printf("Server %s is back to healthy state", server.URL)
			} else {
				log.Printf("Server %s is down", server.URL)
			}
			healthState[server.URL] = server.Health
		}
	}
}
