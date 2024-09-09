package main

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

type VirtualService struct {
	Port       int
	Algorithm  string
	ServerList []*Server
	Logger     *logrus.Logger
}

func (vs *VirtualService) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		forwardRequest(res, req, vs)
	})
	vs.Logger.WithFields(logrus.Fields{
		"port":      vs.Port,
		"algorithm": vs.Algorithm,
	}).Info("Starting virtual service")
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", vs.Port),
		Handler: mux,
	}
	vs.Logger.Fatal(server.ListenAndServe())
}

func forwardRequest(res http.ResponseWriter, req *http.Request, vs *VirtualService) {
	server, err := getHealthyServer(vs)
	if err != nil {
		http.Error(res, "Couldn't process request: "+err.Error(), http.StatusServiceUnavailable)
		return
	}
	server.Connections++
	defer func() { server.Connections-- }()
	vs.Logger.WithFields(logrus.Fields{
		"server": server.URL,
	}).Info("Forwarding request")
	server.ReverseProxy.ServeHTTP(res, req)
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
