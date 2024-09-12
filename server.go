package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Server struct {
	Name         string `json:"name"`
	URL          string `json:"url"`
	Weight       int    `json:"weight"`
	Health       bool   `json:"health"`
	Connections  int    `json:"connections"`
	ReverseProxy *httputil.ReverseProxy
}

func NewServer(name, urlStr string, weight int) *Server {
	url, _ := url.Parse(urlStr)
	return &Server{
		Name:         name,
		URL:          urlStr,
		Weight:       weight,
		Health:       true,
		Connections:  0,
		ReverseProxy: httputil.NewSingleHostReverseProxy(url),
	}
}

func (s *Server) checkHealth() bool {
	resp, err := http.Head(s.URL)
	if err != nil {
		s.Health = false
		return s.Health
	}
	if resp.StatusCode != http.StatusOK {
		s.Health = false
		return s.Health
	}
	s.Health = true
	return s.Health
}
