package registration

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"time"
)

type Registration struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// responsible for registering the service
func RegisterService(webSocketURL, serviceID, serviceType string) error {
	registrationData := Registration{
		Name: serviceID,
		Type: serviceType,
	}

	jsonData, err := json.Marshal(registrationData)
	if err != nil {
		return fmt.Errorf("error marshalling registration data: %w", err)
	}

	err = registerWithRegistry(webSocketURL, jsonData)
	if err != nil {
		return fmt.Errorf("registration failed: %w", err)
	}

	return nil
}

func registerWithRegistry(registryURL string, jsonData []byte) error {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c, _, err := websocket.DefaultDialer.Dial(registryURL, nil)
			if err != nil {
				fmt.Println("Error connecting to WebSocket, retrying...:", err)
				continue
			}

			err = c.WriteMessage(websocket.TextMessage, jsonData)
			if err != nil {
				fmt.Println("Error sending registration data, retrying...:", err)
				c.Close()
				continue
			}

			_, message, err := c.ReadMessage()
			if err != nil {
				fmt.Println("Error reading response, retrying...:", err)
				c.Close()
				continue
			} else {
				fmt.Printf("Response from server: %s\n", message)
			}
			// Consider keeping the connection open if needed for ongoing communication
			// For now, we close it but don't exit the function
			c.Close()
		}
	}
}
