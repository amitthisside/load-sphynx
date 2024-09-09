package main

import (
	"time"

	"github.com/go-co-op/gocron"
)

func startHealthCheck(virtualServices []*VirtualService) {
	for _, vs := range virtualServices {
		logger := vs.Logger
		s := gocron.NewScheduler(time.Local)
		for _, host := range vs.ServerList {
			_, err := s.Every(2).Seconds().Do(func(s *Server) {
				healthy := s.checkHealth()
				if healthy {
					logger.Infof("'%s' is healthy!", s.Name)
				} else {
					logger.Warnf("'%s' is not healthy", s.Name)
				}
			}, host)
			if err != nil {
				logger.Fatalln(err)
			}
		}
		s.StartAsync()
	}
}
