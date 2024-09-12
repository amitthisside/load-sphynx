package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/sirupsen/logrus"
)

type VirtualService struct {
	Port       int            `json:"port"`
	Algorithm  string         `json:"algorithm"`
	ServerList []*Server      `json:"server_list"`
	Logger     *logrus.Logger `json:"-"`
}

func (vs *VirtualService) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		forwardRequest(res, req, vs)
	})
	vs.Logger.Infof("Starting virtual service on port %d with algorithm %s", vs.Port, vs.Algorithm)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", vs.Port),
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}

func forwardRequest(res http.ResponseWriter, req *http.Request, vs *VirtualService) {
	server, err := getHealthyServer(vs)
	if err != nil {
		http.Error(res, "Couldn't process request: "+err.Error(), http.StatusServiceUnavailable)
		return
	}
	server.Connections++
	defer func() { server.Connections-- }()

	vs.Logger.Infof("Forwarding request to: %s", server.URL)
	proxyURL, err := url.Parse(server.URL)
	if err != nil {
		http.Error(res, "Invalid server URL: "+err.Error(), http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(proxyURL)
	proxy.ServeHTTP(res, req)
}

func getHealthyServer(vs *VirtualService) (*Server, error) {
	switch vs.Algorithm {
	case "round_robin":
		return getRoundRobinServer(vs.ServerList, vs.Logger)
	case "weighted_round_robin":
		return getWeightedRoundRobinServer(vs.ServerList, vs.Logger)
	case "least_connections":
		return getLeastConnectionsServer(vs.ServerList, vs.Logger)
	case "weighted_least_connections":
		return getWeightedLeastConnectionsServer(vs.ServerList, vs.Logger)
	default:
		return nil, fmt.Errorf("unknown load balancing algorithm")
	}
}
