package main

import (
	pb "back-end/user"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"log"
	"net/http"
	"os"
)

//type userInput struct {
//	Name          string `json:"name"`
//	Age           string `json:"age"`
//	CommuteMethod string `json:"commute_method"`
//	College       string `json:"college"`
//	Hobbies       string `json:"hobbies"`
//}

type userInput = pb.User

//type User struct {
//	Name          string `protobuf:"bytes,1,opt,name=name,proto3" json:"name"`
//	Age           string `protobuf:"bytes,2,opt,name=age,proto3" json:"age"`
//	CommuteMethod string `protobuf:"bytes,3,opt,name=commute_method,json=commuteMethod,proto3" json:"commute_method"`
//	College       string `protobuf:"bytes,4,opt,name=college,proto3" json:"college"`
//	Hobbies       string `protobuf:"bytes,5,opt,name=hobbies,proto3" json:"hobbies"`
//}

//func (u *User) ProtoReflect() protoreflect.Message {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (u *User) ToProto() ([]byte, error) {
//	return proto.Marshal(u)
//}
//
//func FromProto(data []byte) (*User, error) {
//	user := &User{}
//	err := proto.Unmarshal(data, user)
//	if err != nil {
//		return nil, err
//	}
//	return user, nil
//}

func getUsersFromSheets(c *gin.Context) {
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
	readRange := "Sheet1"

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	var users []userInput
	for _, row := range resp.Values {
		if len(row) >= 5 { // Ensure that there are enough columns in the row
			user := userInput{
				Name:          row[0].(string),
				Age:           row[1].(string), // Assuming age is a number in the sheet
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
	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unable to read client secret file"})
		return
	}

	// If modifying these scopes, delete your previously saved token.json.
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
	var user userInput
	found := false

	for _, row := range resp.Values {
		if len(row) >= 5 {
			if row[0].(string) == name {
				user = userInput{
					Name:          row[0].(string),
					Age:           row[1].(string), // Assuming age is a string in the sheet
					CommuteMethod: row[2].(string),
					College:       row[3].(string),
					Hobbies:       row[4].(string),
				}
				found = true
				break
			}
		}
	}

	if found {
		c.IndentedJSON(http.StatusOK, user)
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
	}
}

func deleteUserFromSheets(c *gin.Context) {
	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unable to read client secret file"})
		return
	}

	// If modifying these scopes, delete your previously saved token.json.
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
	var userRowIndex int
	found := false

	for rowIndex, row := range resp.Values {
		if len(row) >= 5 {
			if row[0].(string) == name {
				userRowIndex = rowIndex
				found = true
				break
			}
		}
	}

	rowToDelete := userRowIndex + 1

	if found {
		deleteRequest := sheets.DeleteDimensionRequest{
			Range: &sheets.DimensionRange{
				SheetId:    getSheetID(spreadsheetID, "Sheet1", srv),
				Dimension:  "ROWS",
				StartIndex: int64(rowToDelete) - 1, // Subtract 1 because row indices are 0-based.
				EndIndex:   int64(rowToDelete),
			},
		}

		// Execute the request to delete the row.
		_, err = srv.Spreadsheets.BatchUpdate(spreadsheetID, &sheets.BatchUpdateSpreadsheetRequest{
			Requests: []*sheets.Request{
				&sheets.Request{
					DeleteDimension: &deleteRequest,
				},
			},
		}).Context(ctx).Do()

		if err != nil {
			log.Fatalf("Error deleting row: %v", err)
		}

		fmt.Printf("Row %d deleted successfully.\n", rowToDelete)

		c.IndentedJSON(http.StatusOK, gin.H{"message": "user deleted"})
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
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
	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unable to read client secret file"})
		return
	}

	// If modifying these scopes, delete your previously saved token.json.
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
	var userRowIndex int
	found := false

	for rowIndex, row := range resp.Values {
		if len(row) >= 5 {
			if row[0].(string) == name {
				userRowIndex = rowIndex
				found = true
				break
			}
		}
	}

	if !found {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
		return
	}

	// Extract the updated user data from the request body
	var updatedUser userInput
	if err := c.BindJSON(&updatedUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid JSON data"})
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
		log.Printf("Error updating user in Google Sheets: %v", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unable to update user in Google Sheets"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "user updated"})
}

func addUser(context *gin.Context) {
	var newUser userInput

	if err := context.BindJSON(&newUser); err != nil {
		return
	}
	if err := addUsertoGoogleSheets(newUser); err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to store data in Google Sheets"})
		return
	}

	context.IndentedJSON(http.StatusCreated, newUser)
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

func main() {
	//initialize the router using gin
	router := gin.Default()

	router.GET("/getuser", getUsersFromSheets)
	router.GET("/getuser/:name", getUserFromSheetsbyName)
	router.PUT("/updateuser/:name", updateUserInSheets)
	router.POST("/deleteuser/:name", deleteUserFromSheets)
	router.POST("/adduser", addUser)
	router.Run("localhost:4000")

}
