package main

var (
	virtualServices = []*VirtualService{
		{
			Port:      8000,
			Algorithm: "weighted_round_robin",
			ServerList: []*Server{
				NewServer("server-1", "http://localhost:5001", 1000),
				NewServer("server-2", "http://localhost:5002", 500),
			},
		},
		{
			Port:      8001,
			Algorithm: "least_connections",
			ServerList: []*Server{
				NewServer("server-3", "http://localhost:5003", 1000),
				NewServer("server-4", "http://localhost:5004", 500),
			},
		},
	}
	healthState = make(map[string]bool)
)

func main() {
	// Initialize logging for each virtual service
	for _, vs := range virtualServices {
		vs.Logger = initLogging(vs)
	}

	// Start health checks
	startHealthCheck(virtualServices)

	for _, vs := range virtualServices {
		go vs.Start()
	}

	select {} // Block forever
}
