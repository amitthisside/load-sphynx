package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

// serverList is a list of URLs of the servers to which the requests will be forwarded.
var (
	serverList      []string
	lastServedIndex = 0
)

func main() {
	loadConfig("config.json")
	http.HandleFunc("/", forwardResponse)
	log.Fatal(http.ListenAndServe(":8000", nil))
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

// to get the next server url to forward the request to
func getServerUrl() *url.URL {
	nextIndex := (lastServedIndex + 1) % len(serverList)
	url, err := url.Parse(serverList[nextIndex])
	lastServedIndex = nextIndex
	if err != nil {
		log.Fatal(err)
	}
	return url
}
