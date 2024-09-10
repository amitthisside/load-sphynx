package main

var (
	virtualServices []*VirtualService
)

func main() {
	// Initialize configuration
	initConfig()

	// Display loaded configuration
	displayConfig()

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
