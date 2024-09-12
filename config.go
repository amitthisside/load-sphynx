package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func initConfig() {
	configFile := "config.json"
	absConfigFile, err := filepath.Abs(configFile)
	if err != nil {
		fmt.Printf("Error getting absolute path of config file: %v\n", err)
		os.Exit(1)
	}

	if _, err := os.Stat(absConfigFile); os.IsNotExist(err) {
		fmt.Println("No config file found, entering configuration manually.")
		readConfigFromInput()
	} else {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Config file found. Do you want to use it? (y/n): ")
		useConfig, _ := reader.ReadString('\n')
		useConfig = strings.TrimSpace(useConfig)
		if useConfig == "y" {
			data, err := os.ReadFile(absConfigFile)
			if err != nil {
				fmt.Printf("Error reading config file: %v\n", err)
				os.Exit(1)
			}
			err = json.Unmarshal(data, &virtualServices)
			if err != nil {
				fmt.Printf("Error parsing config file: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Config file successfully loaded.")
		} else {
			fmt.Println("Entering configuration manually.")
			readConfigFromInput()
		}
	}
}

func readConfigFromInput() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter port: ")
		portStr, _ := reader.ReadString('\n')
		port, _ := strconv.Atoi(strings.TrimSpace(portStr))

		fmt.Print("Enter algorithm: ")
		algorithm, _ := reader.ReadString('\n')
		algorithm = strings.TrimSpace(algorithm)

		var serverList []*Server
		for {
			fmt.Print("Enter server name (or 'done' to finish): ")
			name, _ := reader.ReadString('\n')
			name = strings.TrimSpace(name)
			if name == "done" {
				break
			}

			fmt.Print("Enter server URL: ")
			url, _ := reader.ReadString('\n')
			url = strings.TrimSpace(url)

			fmt.Print("Enter server weight: ")
			weightStr, _ := reader.ReadString('\n')
			weight, _ := strconv.Atoi(strings.TrimSpace(weightStr))

			serverList = append(serverList, NewServer(name, url, weight))
		}

		// Ensure we append new VirtualService to the global slice
		virtualServices = append(virtualServices, &VirtualService{
			Port:       port,
			Algorithm:  algorithm,
			ServerList: serverList,
			Logger:     initLogging(&VirtualService{Port: port, Algorithm: algorithm}),
		})

		fmt.Print("Add another virtual service? (y/n): ")
		another, _ := reader.ReadString('\n')
		another = strings.TrimSpace(another)
		if another != "y" {
			break
		}
	}
}

func displayConfig() {
	configBytes, err := json.MarshalIndent(virtualServices, "", "  ")
	if err != nil {
		fmt.Printf("Error displaying config: %v\n", err)
		return
	}
	fmt.Printf("Loaded configuration:\n%s\n", string(configBytes))
}
