package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"
)

var (
	serverList      []string
	lastServedIndex = 0
	mu              sync.Mutex // Mutex to protect lastServedIndex
)

func main() {
	loadConfig("server_conf.json")
	http.HandleFunc("/", forwardResponse)
	port := ":8000"
	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// loadConfig reads the server list from a JSON configuration file.
func loadConfig(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := struct {
		Servers []string `json:"servers"`
	}{}

	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	serverList = config.Servers
}

// forwardResponse forwards the HTTP request to the server specified by the URL returned from getServerUrl().
// It uses a reverse proxy to handle the forwarding process.
// The response from the server is written to the provided http.ResponseWriter.
// The request details are passed through the req parameter.
func forwardResponse(res http.ResponseWriter, req *http.Request) {
	url := getServerUrl()
	rProxy := httputil.NewSingleHostReverseProxy(url)
	log.Printf("Forwarding request to %s", url.String())
	rProxy.ServeHTTP(res, req)
}

// getServerUrl returns the next server URL to forward the request to.
func getServerUrl() *url.URL {
	mu.Lock()
	defer mu.Unlock()

	if len(serverList) == 0 {
		log.Fatal("No servers available")
	}

	nextIndex := (lastServedIndex + 1) % len(serverList)
	serverURL, err := url.Parse(serverList[nextIndex])
	if err != nil {
		log.Fatalf("Failed to parse server URL: %v", err)
	}

	lastServedIndex = nextIndex
	return serverURL
}
