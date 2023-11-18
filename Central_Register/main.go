package main

import (
	gs "centralReg/grpc_status"
	pb "centralReg/service_reg"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"os"
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

	// Use the PORT environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8090" // Default to 8090 if PORT is not set
	}

	http.HandleFunc("/register", registerService)
	http.HandleFunc("/services", listServices)

	fmt.Printf("Service Registry is running on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func monitorServices() {
	for {
		registry.RLock()
		for _, service := range registry.services {
			//go checkServiceStatus(service)
			if service.Type == "gRPC" {
				go checkGRPCServiceStatus(service)
			} else {
				go checkRESTServiceStatus(service)
			}

		}
		registry.RUnlock()
		time.Sleep(10 * time.Second) // Check every 10 seconds
	}
}

func checkRESTServiceStatus(service *ApiService) {
	var url string
	if service.Host == "localhost" {
		url = fmt.Sprintf("http://%s:%d/status", service.Host, service.Port)
	} else {
		url = fmt.Sprintf("https://%s/status", service.Host)
	}

	// Print the URL to the console
	fmt.Println("Checking service status at:", url)

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		service.Status = "Down"
	} else {
		service.Status = "Up"
	}

	// It's important to close the response body to free resources
	if resp != nil {
		resp.Body.Close()
	}
}

func checkGRPCServiceStatus(service *ApiService) {
	// Set up a connection to the server with insecure credentials for simplicity
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", service.Host, service.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to dial gRPC service: %v", err)
		service.Status = "Down"
		return
	}
	defer conn.Close()

	// Create a new StatusService client
	client := gs.NewStatusServiceClient(conn)

	// Prepare a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Call the CheckStatus method
	response, err := client.CheckStatus(ctx, &gs.StatusRequest{})
	if err != nil {
		log.Printf("Error calling CheckStatus: %v", err)
		service.Status = "Down"
	} else {
		service.Status = response.Status
		//service.Status = "Up"
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
			conn.WriteMessage(websocket.TextMessage, []byte("Invalid registration data"))
			continue
		}

		registry.Lock()
		if existingService, exists := registry.services[reg.Name]; exists {
			// Optionally update existing registration instead of ignoring
			existingService.Host = reg.Host
			existingService.Port = reg.Port
			existingService.Type = reg.Type
			log.Printf("Updated registration for service %s", reg.Name)
		} else {
			registry.services[reg.Name] = &ApiService{
				Name: reg.Name,
				Host: reg.Host,
				Port: reg.Port,
				Type: reg.Type,
			}
			log.Printf("Service %s registered successfully with Host: %s, Port: %d", reg.Name, reg.Host, reg.Port)
		}
		registry.Unlock()

		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Service %s processed", reg.Name)))
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
