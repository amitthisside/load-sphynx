package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Get all virtual services
func getVirtualServices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(virtualServices); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// Get a specific virtual service by Port (used as ID here)
func getVirtualService(w http.ResponseWriter, r *http.Request) {
	portStr := mux.Vars(r)["vs_id"]
	port, err := strconv.Atoi(portStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	for _, vs := range virtualServices {
		if vs.Port == port {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(vs)
			return
		}
	}
	http.Error(w, "Virtual service not found", http.StatusNotFound)
}

// Create a new virtual service
func createVirtualService(w http.ResponseWriter, r *http.Request) {
	var vs VirtualService
	if err := json.NewDecoder(r.Body).Decode(&vs); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	vs.Logger = initLogging(&vs) // Initialize logger for the new VS
	virtualServices = append(virtualServices, &vs)
	w.WriteHeader(http.StatusCreated)
}

// Update a virtual service
func updateVirtualService(w http.ResponseWriter, r *http.Request) {
	portStr := mux.Vars(r)["vs_id"]
	port, err := strconv.Atoi(portStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var updatedVS *VirtualService
	err = json.NewDecoder(r.Body).Decode(&updatedVS)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	for i, vs := range virtualServices {
		if vs.Port == port {
			virtualServices[i] = updatedVS
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	http.Error(w, "Virtual service not found", http.StatusNotFound)
}

// Delete a virtual service
func deleteVirtualService(w http.ResponseWriter, r *http.Request) {
	portStr := mux.Vars(r)["vs_id"]
	port, err := strconv.Atoi(portStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	for i, vs := range virtualServices {
		if vs.Port == port {
			virtualServices = append(virtualServices[:i], virtualServices[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Virtual service not found", http.StatusNotFound)
}
