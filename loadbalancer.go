package main

import (
	"fmt"
	"log"
	"net/http"
)

var (
	serverList = []*server{
		newServer("server-1", "http://localhost:5001"),
		newServer("server-2", "http://localhost:5002"),
		newServer("server-3", "http://localhost:5003"),
		newServer("server-4", "http://localhost:5004"),
		newServer("server-5", "http://localhost:5005"),
	}
	lastServedIndex = 0
)

func main() {
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
	server.ReverseProxy.ServeHTTP(res, req)
}

func getHealthyServer() (*server, error) {
	for i := 0; i < len(serverList); i++ {
		server := getServer()
		if server.Health {
			return server, nil
		}
	}
	return nil, fmt.Errorf("no healthy hosts")
}

func getServer() *server {
	nextIndex := (lastServedIndex + 1) % len(serverList)
	server := serverList[nextIndex]
	lastServedIndex = nextIndex
	return server
}
