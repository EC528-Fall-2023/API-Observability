package main

import (
	"back-end/influxdb"
	pb "back-end/user"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"log"
	"net/http"
	"os"
	"time"
)

type userInput = pb.User

var influxClient = influxdb.InitInfluxDB()

var getUserCount int = 0
var getUserError int = 0

var getUserNameCount int = 0
var getUserNameError int = 0

var updateUserCount int = 0
var updateUserError int = 0

var deleteUserCount int = 0
var deleteUserError int = 0

var addUserCount int = 0
var addUserError int = 0

func init() {
	fmt.Println("init called")
	influxClient = influxdb.InitInfluxDB()
}

func getUsersFromSheets(c *gin.Context) {
	startTime := time.Now()
	getUserCount++
	tags := map[string]string{
		"endpoint":   "getuser",
		"user_agent": c.Request.UserAgent(),
		"ip_address": c.ClientIP(),
	}
	fields := map[string]interface{}{
		"request_count": int64(getUserCount),
		"request_size":  c.Request.ContentLength,
	}

	var users []userInput
	defer func() {
		if r := recover(); r != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}

		if len(users) == 0 {
			getUserError++
			fields["error_count"] = int64(getUserError)
			fields["status_code"] = http.StatusInternalServerError
		} else {
			fields["status_code"] = http.StatusOK
		}

		fields["response_size"] = c.Writer.Size()

		latency := time.Since(startTime)
		latencyMs := float64(latency.Nanoseconds()) / 1000000.0
		fields["latency"] = latencyMs

		if err := influxdb.WriteMetric(influxClient, "REST", tags, fields); err != nil {
			fmt.Println("Error writing metrics:", err)
		}
	}()

	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	spreadsheetId := "10-CfbfktbeTSMV3tgnIKwaBquzw-RmjS13Tut9A32_s"
	readRange := "Sheet1"

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	for _, row := range resp.Values {
		if len(row) >= 5 {
			user := userInput{
				Name:          row[0].(string),
				Age:           row[1].(string),
				CommuteMethod: row[2].(string),
				College:       row[3].(string),
				Hobbies:       row[4].(string),
			}
			users = append(users, user)
		}
	}
	c.IndentedJSON(http.StatusOK, users)
}

func getUserFromSheetsbyName(c *gin.Context) {
	startTime := time.Now()
	getUserNameCount++
	tags := map[string]string{
		"endpoint":   "getuser/name",
		"user_agent": c.Request.UserAgent(),
		"ip_address": c.ClientIP(),
	}
	fields := map[string]interface{}{
		"request_count": int64(getUserNameCount),
		"request_size":  c.Request.ContentLength,
	}

	var user userInput
	defer func() {
		if r := recover(); r != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}

		if user.Name == "" {
			getUserNameError++
			fields["error_count"] = int64(getUserNameError)
			fields["status_code"] = http.StatusNotFound
		} else {
			fields["status_code"] = http.StatusOK
		}

		fields["response_size"] = c.Writer.Size()

		latency := time.Since(startTime)
		latencyMs := float64(latency.Nanoseconds()) / 1000000.0
		fields["latency"] = latencyMs

		if err := influxdb.WriteMetric(influxClient, "REST", tags, fields); err != nil {
			fmt.Println("Error writing metrics:", err)
		}
	}()

	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unable to read client secret file"})
		return
	}

	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unable to parse client secret file to config"})
		return
	}
	client := getClient(config)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unable to retrieve Sheets client"})
		return
	}

	spreadsheetID := "10-CfbfktbeTSMV3tgnIKwaBquzw-RmjS13Tut9A32_s"
	readRange := "Sheet1"

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unable to retrieve data from Google Sheets"})
		return
	}

	name := c.Param("name")
	found := false

	for _, row := range resp.Values {
		if len(row) >= 5 && row[0].(string) == name {
			user = userInput{
				Name:          row[0].(string),
				Age:           row[1].(string),
				CommuteMethod: row[2].(string),
				College:       row[3].(string),
				Hobbies:       row[4].(string),
			}
			found = true
			break
		}
	}

	if found {
		c.IndentedJSON(http.StatusOK, user)
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "User not found"})
	}
}

func deleteUserFromSheets(c *gin.Context) {
	startTime := time.Now()
	deleteUserCount++
	tags := map[string]string{
		"endpoint":   "deleteuser/name",
		"user_agent": c.Request.UserAgent(),
		"ip_address": c.ClientIP(),
	}
	fields := map[string]interface{}{
		"request_count": int64(deleteUserCount),
		"request_size":  c.Request.ContentLength,
	}

	var status int
	defer func() {
		fields["response_size"] = c.Writer.Size()
		fields["status_code"] = status

		latency := time.Since(startTime)
		latencyMs := float64(latency.Nanoseconds()) / 1000000.0
		fields["latency"] = latencyMs

		if err := influxdb.WriteMetric(influxClient, "REST", tags, fields); err != nil {
			fmt.Println("Error writing metrics:", err)
		}
	}()

	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		status = http.StatusInternalServerError
		c.IndentedJSON(status, gin.H{"message": "Unable to read client secret file"})
		return
	}

	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		status = http.StatusInternalServerError
		c.IndentedJSON(status, gin.H{"message": "Unable to parse client secret file to config"})
		return
	}
	client := getClient(config)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		status = http.StatusInternalServerError
		c.IndentedJSON(status, gin.H{"message": "Unable to retrieve Sheets client"})
		return
	}

	spreadsheetID := "10-CfbfktbeTSMV3tgnIKwaBquzw-RmjS13Tut9A32_s"
	readRange := "Sheet1"

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		status = http.StatusInternalServerError
		c.IndentedJSON(status, gin.H{"message": "Unable to retrieve data from Google Sheets"})
		return
	}

	name := c.Param("name")
	var userRowIndex int
	found := false

	for rowIndex, row := range resp.Values {
		if len(row) >= 5 && row[0].(string) == name {
			userRowIndex = rowIndex
			found = true
			break
		}
	}

	if found {
		rowToDelete := userRowIndex + 1
		deleteRequest := sheets.DeleteDimensionRequest{
			Range: &sheets.DimensionRange{
				SheetId:    getSheetID(spreadsheetID, "Sheet1", srv),
				Dimension:  "ROWS",
				StartIndex: int64(rowToDelete) - 1, // Subtract 1 because row indices are 0-based.
				EndIndex:   int64(rowToDelete),
			},
		}

		_, err = srv.Spreadsheets.BatchUpdate(spreadsheetID, &sheets.BatchUpdateSpreadsheetRequest{
			Requests: []*sheets.Request{
				{DeleteDimension: &deleteRequest},
			},
		}).Context(ctx).Do()

		if err != nil {
			deleteUserError++
			fields["error_count"] = int64(deleteUserError)
			status = http.StatusInternalServerError
			log.Printf("Error deleting row: %v", err)
			c.IndentedJSON(status, gin.H{"message": "Error deleting user"})
			return
		}

		status = http.StatusOK
		fmt.Printf("Row %d deleted successfully.\n", rowToDelete)
		c.IndentedJSON(status, gin.H{"message": "User deleted"})
	} else {
		status = http.StatusNotFound
		c.IndentedJSON(status, gin.H{"message": "User not found"})
	}
}

func getSheetID(spreadsheetID, sheetName string, srv *sheets.Service) int64 {
	resp, err := srv.Spreadsheets.Get(spreadsheetID).Context(context.Background()).Do()
	if err != nil {
		log.Printf("Unable to retrieve spreadsheet: %v", err)
		return -1
	}

	for _, sheet := range resp.Sheets {
		if sheet.Properties.Title == sheetName {
			return sheet.Properties.SheetId
		}
	}

	log.Printf("Sheet not found: %s", sheetName)
	return -1
}

func updateUserInSheets(c *gin.Context) {
	startTime := time.Now()
	updateUserCount++
	tags := map[string]string{
		"endpoint":   "updateuser/name",
		"user_agent": c.Request.UserAgent(),
		"ip_address": c.ClientIP(),
	}
	fields := map[string]interface{}{
		"request_count": int64(updateUserCount),
		"request_size":  c.Request.ContentLength,
	}

	var status int
	defer func() {
		fields["response_size"] = c.Writer.Size()
		fields["status_code"] = status

		latency := time.Since(startTime)
		latencyMs := float64(latency.Nanoseconds()) / 1000000.0
		fields["latency"] = latencyMs

		if err := influxdb.WriteMetric(influxClient, "REST", tags, fields); err != nil {
			fmt.Println("Error writing metrics:", err)
		}
	}()

	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		status = http.StatusInternalServerError
		c.IndentedJSON(status, gin.H{"message": "Unable to read client secret file"})
		return
	}

	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		status = http.StatusInternalServerError
		c.IndentedJSON(status, gin.H{"message": "Unable to parse client secret file to config"})
		return
	}
	client := getClient(config)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		status = http.StatusInternalServerError
		c.IndentedJSON(status, gin.H{"message": "Unable to retrieve Sheets client"})
		return
	}

	spreadsheetID := "10-CfbfktbeTSMV3tgnIKwaBquzw-RmjS13Tut9A32_s"
	readRange := "Sheet1"

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		status = http.StatusInternalServerError
		c.IndentedJSON(status, gin.H{"message": "Unable to retrieve data from Google Sheets"})
		return
	}

	name := c.Param("name")
	var userRowIndex int
	found := false

	for rowIndex, row := range resp.Values {
		if len(row) >= 5 && row[0].(string) == name {
			userRowIndex = rowIndex
			found = true
			break
		}
	}

	if !found {
		status = http.StatusNotFound
		c.IndentedJSON(status, gin.H{"message": "User not found"})
		return
	}

	// Extract the updated user data from the request body
	var updatedUser userInput
	if err := c.BindJSON(&updatedUser); err != nil {
		status = http.StatusBadRequest
		c.IndentedJSON(status, gin.H{"message": "Invalid JSON data"})
		return
	}

	// Update the user's data in the spreadsheet
	writeRange := fmt.Sprintf("Sheet1!A%d:E%d", userRowIndex+1, userRowIndex+1)
	var values [][]interface{}
	values = append(values, []interface{}{updatedUser.Name, updatedUser.Age, updatedUser.CommuteMethod, updatedUser.College, updatedUser.Hobbies})

	vr := &sheets.ValueRange{
		Values: values,
	}

	_, err = srv.Spreadsheets.Values.Update(spreadsheetID, writeRange, vr).ValueInputOption("RAW").Do()
	if err != nil {
		updateUserError++
		fields["error_count"] = int64(updateUserError)
		status = http.StatusInternalServerError
		log.Printf("Error updating user in Google Sheets: %v", err)
		c.IndentedJSON(status, gin.H{"message": "Unable to update user in Google Sheets"})
		return
	}

	status = http.StatusOK
	c.IndentedJSON(status, gin.H{"message": "User updated"})
}

func addUser(context *gin.Context) {
	startTime := time.Now()
	addUserCount++
	tags := map[string]string{
		"endpoint":   "adduser",
		"user_agent": context.Request.UserAgent(),
		"ip_address": context.ClientIP(),
	}
	fields := map[string]interface{}{
		"request_count": int64(addUserCount),
		"request_size":  context.Request.ContentLength,
	}

	var newUser userInput
	var status int

	defer func() {
		fields["response_size"] = context.Writer.Size()
		fields["status_code"] = status

		latency := time.Since(startTime)
		latencyMs := float64(latency.Nanoseconds()) / 1000000.0
		fields["latency"] = latencyMs

		if err := influxdb.WriteMetric(influxClient, "REST", tags, fields); err != nil {
			fmt.Println("Error writing metrics:", err)
		}
	}()

	if err := context.BindJSON(&newUser); err != nil {
		status = http.StatusBadRequest
		context.IndentedJSON(status, gin.H{"message": "Invalid JSON data"})
		return
	}
	if err := addUsertoGoogleSheets(newUser); err != nil {
		addUserError++
		fields["error_count"] = int64(addUserError)
		status = http.StatusInternalServerError
		context.IndentedJSON(status, gin.H{"message": "Failed to store data in Google Sheets"})
		return
	}

	status = http.StatusCreated
	context.IndentedJSON(status, newUser)
}

// These functions are used for the google sheets API
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// functions for interacting with google sheets
func addUsertoGoogleSheets(user userInput) error {
	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	spreadsheetId := "10-CfbfktbeTSMV3tgnIKwaBquzw-RmjS13Tut9A32_s"
	writeRange := "Sheet1"

	var values [][]interface{}
	values = append(values, []interface{}{user.Name, user.Age, user.CommuteMethod, user.College, user.Hobbies})

	vr := &sheets.ValueRange{
		Values: values,
	}

	_, err = srv.Spreadsheets.Values.Append(spreadsheetId, writeRange, vr).ValueInputOption("RAW").Do()
	if err != nil {
		return err
	}

	return nil
}

func registerWithRegistry(name, host string, port int, servType string) {
	registryURL := "ws://localhost:8090/register" // WebSocket URL
	registrationData := Registration{
		Name: name,
		Host: host,
		Port: port,
		Type: servType,
	}

	jsonData, err := json.Marshal(registrationData)
	if err != nil {
		fmt.Println("Error marshalling registration data:", err)
		return
	}

	ticker := time.NewTicker(10 * time.Second) // Retry every 10 seconds
	defer ticker.Stop()

	for {
		c, _, err := websocket.DefaultDialer.Dial(registryURL, nil)
		if err != nil {
			fmt.Println("Error connecting to WebSocket, retrying...:", err)
		} else {
			// Successfully connected, send registration data
			err = c.WriteMessage(websocket.TextMessage, jsonData)
			if err != nil {
				fmt.Println("Error sending registration data, retrying...:", err)
			} else {
				// Read response
				_, message, err := c.ReadMessage()
				if err != nil {
					fmt.Println("Error reading response, retrying...:", err)
				} else {
					fmt.Printf("Response from server: %s\n", message)
					c.Close()
					break // Exit the loop if registration is successful
				}
			}
			c.Close() // Close the connection in case of any error
		}

		// Wait for the next tick before retrying
		<-ticker.C
	}
}

type Registration struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Port int    `json:"port"`
	Type string `json:"type"`
}

func main() {
	serviceName := "Student-Info REST Service"
	serviceHost := "localhost"
	servicePort := 4000
	serviceType := "REST"

	// Register your service with the registry
	go registerWithRegistry(serviceName, serviceHost, servicePort, serviceType)

	//initialize the router using gin
	router := gin.Default()

	router.GET("/getuser", getUsersFromSheets)
	router.GET("/getuser/:name", getUserFromSheetsbyName)
	router.PUT("/updateuser/:name", updateUserInSheets)
	router.POST("/deleteuser/:name", deleteUserFromSheets)
	router.POST("/adduser", addUser)
	router.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "Up",
		})
	})
	router.Run("localhost:4000")

}
