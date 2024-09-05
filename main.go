package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// serverList is a list of URLs of the servers to which the requests will be forwarded.
var (
	serverList = []string{
		"http://localhost:5000",
		"http://localhost:5001",
		"http://localhost:5002",
		"http://localhost:5003",
		"http://localhost:5004",
	}
	lastServedIndex = 0
)

func main() {
	http.HandleFunc("/", forwardResponse)
	log.Fatal(http.ListenAndServe(":8000", nil))
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
