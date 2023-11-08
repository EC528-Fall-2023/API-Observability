// service-registry.go
package main

import (
	pb "centralReg/service_reg"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type ApiService = pb.APIService
type Registration = pb.Registration

var registry = struct {
	sync.RWMutex
	services map[string]*ApiService
}{services: make(map[string]*ApiService)}

func main() {
	go monitorServices()
	http.HandleFunc("/register", registerService)
	http.HandleFunc("/services", listServices)
	fmt.Println("Service Registry is running on port 8090")
	http.ListenAndServe(":8090", nil)
}

func monitorServices() {
	for {
		registry.RLock()
		for _, service := range registry.services {
			go checkServiceStatus(service)
		}
		registry.RUnlock()
		time.Sleep(10 * time.Second) // Check every 10 seconds
	}
}

func checkServiceStatus(service *ApiService) {
	url := fmt.Sprintf("http://%s:%d/status", service.Host, service.Port)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		service.Status = "Down"
	} else {
		service.Status = "Up"
	}
}

func registerService(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var reg Registration
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	registry.Lock()
	registry.services[reg.Name] = &ApiService{
		Name: reg.Name,
		Host: reg.Host,
		Port: reg.Port,
		Type: reg.Type,
	}
	registry.Unlock()

	log.Printf("Service %s which is %s registered successfully with Host: %s, Port: %d\n", reg.Name, reg.Type, reg.Host, reg.Port)
	fmt.Fprintf(w, "Service %s registered successfully", reg.Name)
}

func listServices(w http.ResponseWriter, r *http.Request) {
	registry.RLock()
	defer registry.RUnlock()

	services := make([]*ApiService, 0, len(registry.services))

	for _, service := range registry.services {
		services = append(services, service)
	}

	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")

	// Encode the services slice to JSON and send it as the response
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(services); err != nil {
		// Handle the error appropriately
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error encoding services list: %s\n", err)
		return
	}

	log.Println("Listing all registered services:")
	// Logging the services; consider removing if it's too verbose
	for _, service := range services {
		log.Printf("%s: %s\n", service.Name, service.Status)
	}
}
