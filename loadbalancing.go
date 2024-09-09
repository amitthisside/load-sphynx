package main

import (
	"fmt"
)

var (
	currentIndex  = -1
	currentWeight = 0
)

func getRoundRobinServer() (*server, error) {
	for i := 0; i < len(serverList); i++ {
		server := serverList[(currentIndex+i)%len(serverList)]
		if server.Health {
			currentIndex = (currentIndex + i + 1) % len(serverList)
			return server, nil
		}
	}
	return nil, fmt.Errorf("no healthy hosts")
}

func getWeightedRoundRobinServer() (*server, error) {
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
			logServerHealthChanges()
			return serverList[currentIndex], nil
		}
	}
}

func getLeastConnectionsServer() (*server, error) {
	var leastConnServer *server
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
	logServerHealthChanges()
	return leastConnServer, nil
}

func getWeightedLeastConnectionsServer() (*server, error) {
	var bestServer *server
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
	logServerHealthChanges()
	return bestServer, nil
}
