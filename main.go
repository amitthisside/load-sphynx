package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	serverList = []*server{
		newServer("server-1", "http://localhost:5001", 1000),
		newServer("server-2", "http://localhost:5002", 500),
		newServer("server-3", "http://localhost:5003", 1000),
		newServer("server-4", "http://localhost:5004", 500),
		newServer("server-5", "http://localhost:5005", 500),
	}
	lbAlgorithm = "weighted_round_robin" // Options: "round_robin", "weighted_round_robin", "least_connections", "weighted_least_connections"
	logFile     *os.File
	healthState = make(map[string]bool)
)

func main() {
	// Initialize logging
	initLogging()

	http.HandleFunc("/", forwardRequest)
	go startHealthCheck()
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func forwardRequest(res http.ResponseWriter, req *http.Request) {
	server, err := getHealthyServer()
	if err != nil {
		http.Error(res, "Couldn't process request: "+err.Error(), http.StatusServiceUnavailable)
		return
	}
	server.Connections++
	defer func() { server.Connections-- }()
	log.Printf("Forwarding request to: %s", server.URL)
	server.ReverseProxy.ServeHTTP(res, req)
}

func getHealthyServer() (*server, error) {
	switch lbAlgorithm {
	case "round_robin":
		return getRoundRobinServer()
	case "weighted_round_robin":
		return getWeightedRoundRobinServer()
	case "least_connections":
		return getLeastConnectionsServer()
	case "weighted_least_connections":
		return getWeightedLeastConnectionsServer()
	default:
		return nil, fmt.Errorf("unknown load balancing algorithm")
	}
}
