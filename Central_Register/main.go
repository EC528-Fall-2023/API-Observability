package main

import (
	pb "centralReg/service_reg"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

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
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade to websocket failed:", err)
		return
	}
	defer conn.Close()

	for {
		// Read message from WebSocket connection
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		var reg Registration
		if err := json.Unmarshal(msg, &reg); err != nil {
			log.Printf("Error decoding registration data: %s", err)
			// Optionally, send an error message back to the client
			conn.WriteMessage(websocket.TextMessage, []byte("Invalid registration data"))
			continue
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
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Service %s registered successfully", reg.Name)))
	}
}

func listServices(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade to websocket failed:", err)
		return
	}
	defer conn.Close()

	registry.RLock()
	defer registry.RUnlock()

	services := make([]*ApiService, 0, len(registry.services))

	for _, service := range registry.services {
		services = append(services, service)
	}

	// Send the service list over WebSocket
	err = conn.WriteJSON(services)
	if err != nil {
		log.Println("Error sending services over WebSocket:", err)
		return
	}

	log.Println("Sent services list over WebSocket")
}
