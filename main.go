package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

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

	// Initialize router
	r := mux.NewRouter()
	r.HandleFunc("/access/vs", getVirtualServices).Methods("GET")
	r.HandleFunc("/access/vs/{vs_id:[0-9]+}", getVirtualService).Methods("GET")
	r.HandleFunc("/access/vs", createVirtualService).Methods("POST")
	r.HandleFunc("/access/vs/{vs_id:[0-9]+}", updateVirtualService).Methods("PUT")
	r.HandleFunc("/access/vs/{vs_id:[0-9]+}", deleteVirtualService).Methods("DELETE")

	// Start HTTP server
	log.Fatal(http.ListenAndServe(":8080", r))

	select {} // Block forever
}
