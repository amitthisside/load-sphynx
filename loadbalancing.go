package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

var (
	currentIndex  = -1
	currentWeight = 0
)

func getRoundRobinServer(serverList []*Server, logger *logrus.Logger) (*Server, error) {
	for i := 0; i < len(serverList); i++ {
		server := serverList[(currentIndex+i)%len(serverList)]
		if server.Health {
			currentIndex = (currentIndex + i + 1) % len(serverList)
			return server, nil
		}
	}
	return nil, fmt.Errorf("no healthy hosts")
}

func getWeightedRoundRobinServer(serverList []*Server, logger *logrus.Logger) (*Server, error) {
	totalWeight := 0
	for _, server := range serverList {
		if server.Health {
			totalWeight += server.Weight
		}
	}
	if totalWeight == 0 {
		return nil, fmt.Errorf("no healthy hosts")
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
		return nil, fmt.Errorf("no healthy hosts")
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
		return nil, fmt.Errorf("no healthy hosts")
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
