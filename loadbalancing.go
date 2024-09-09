package main

import (
	"errors"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	rrIndex       int
	rrMutex       sync.Mutex
	currentIndex  int
	currentWeight int

	healthState = make(map[string]bool)
)

func getRoundRobinServer(serverList []*Server, logger *logrus.Logger) (*Server, error) {
	rrMutex.Lock()
	defer rrMutex.Unlock()

	if len(serverList) == 0 {
		return nil, errors.New("no servers available")
	}

	for {
		rrIndex = (rrIndex + 1) % len(serverList)
		server := serverList[rrIndex]
		if server.Health {
			logServerHealthChanges(serverList, logger)
			return server, nil
		}
	}
}

func getWeightedRoundRobinServer(serverList []*Server, logger *logrus.Logger) (*Server, error) {
	totalWeight := 0
	for _, server := range serverList {
		if server.Health {
			totalWeight += server.Weight
		}
	}
	if totalWeight == 0 {
		return nil, errors.New("no healthy hosts")
	}

	for {
		currentIndex = (currentIndex + 1) % len(serverList)
		if currentIndex == 0 {
			currentWeight = currentWeight - gcd(serverList)
			if currentWeight <= 0 {
				currentWeight = maxWeight(serverList)
			}
		}
		if serverList[currentIndex].Health && serverList[currentIndex].Weight >= currentWeight {
			logServerHealthChanges(serverList, logger)
			return serverList[currentIndex], nil
		}
	}
}

func getLeastConnectionsServer(serverList []*Server, logger *logrus.Logger) (*Server, error) {
	var leastConnServer *Server
	for _, server := range serverList {
		if server.Health {
			if leastConnServer == nil || server.Connections < leastConnServer.Connections {
				leastConnServer = server
			}
		}
	}
	if leastConnServer == nil {
		return nil, errors.New("no healthy hosts")
	}
	logServerHealthChanges(serverList, logger)
	return leastConnServer, nil
}

func getWeightedLeastConnectionsServer(serverList []*Server, logger *logrus.Logger) (*Server, error) {
	var bestServer *Server
	var minRatio float64 = -1
	for _, server := range serverList {
		if server.Health {
			ratio := float64(server.Connections) / float64(server.Weight)
			if minRatio == -1 || ratio < minRatio {
				minRatio = ratio
				bestServer = server
			}
		}
	}
	if bestServer == nil {
		return nil, errors.New("no healthy hosts")
	}
	logServerHealthChanges(serverList, logger)
	return bestServer, nil
}

func logServerHealthChanges(serverList []*Server, logger *logrus.Logger) {
	for _, server := range serverList {
		previousHealth, exists := healthState[server.URL]
		if !exists || previousHealth != server.Health {
			if server.Health {
				logger.WithFields(logrus.Fields{
					"server": server.URL,
				}).Info("Server is back to healthy state")
			} else {
				logger.WithFields(logrus.Fields{
					"server": server.URL,
				}).Warn("Server is down")
			}
			healthState[server.URL] = server.Health
		}
	}
}
