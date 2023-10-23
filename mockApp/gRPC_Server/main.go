// main.go
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/Ling-Qingran/gRPC-Observability/user"
	"google.golang.org/grpc"
)

type userServiceServer struct {
	user.UnimplementedUserServiceServer
}

const (
	spreadsheetID = "10-CfbfktbeTSMV3tgnIKwaBquzw-RmjS13Tut9A32_s"
	readRange     = "Sheet1"
	writeRange    = "Sheet1"
)

var srv *sheets.Service

// Retrieve a token, saves the token, then returns the generated client.
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

func init() {
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

	srv, err = sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}
}

func (s *userServiceServer) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.User, error) {
	name := req.GetName()

	// Get users from Google Sheets
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		return nil, err
	}

	// Iterate through the rows to find the user with the given name
	for _, row := range resp.Values {
		if len(row) < 5 {
			continue
		}
		if row[0].(string) == name {
			return &user.User{
				Name:          row[0].(string),
				Age:           row[1].(string), // assuming age is stored as a number
				CommuteMethod: row[2].(string),
				College:       row[3].(string),
				Hobbies:       row[4].(string),
			}, nil
		}
	}

	return nil, errors.New("user not found")
}

func (s *userServiceServer) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*user.User, error) {
	name := req.GetName()

	rowNumber, err := getRowNumberByName(name)
	if err != nil {
		return nil, err
	}

	updatedUser := req.GetUser()

	// Prepare the data for the update in Google Sheets
	var rowData []interface{}
	rowData = append(rowData, updatedUser.Name)
	rowData = append(rowData, updatedUser.Age)
	rowData = append(rowData, updatedUser.CommuteMethod)
	rowData = append(rowData, updatedUser.College)
	rowData = append(rowData, updatedUser.Hobbies)

	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{rowData},
	}

	updateRange := fmt.Sprintf("Sheet1!A%d:E%d", rowNumber, rowNumber) // Assuming data starts in column A and spans 5 columns
	_, err = srv.Spreadsheets.Values.Update(spreadsheetID, updateRange, valueRange).ValueInputOption("RAW").Do()
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func getRowNumberByName(name string) (int, error) {
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do() // Assuming name is in column A
	if err != nil {
		return -1, err
	}

	for i, row := range resp.Values {
		if len(row) > 0 && row[0].(string) == name {
			return i + 1, nil // +1 because sheet row numbers start at 1, not 0
		}
	}

	return -1, errors.New("user not found")
}

func (s *userServiceServer) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (*user.DeleteUserResponse, error) {
	name := req.GetName()

	rowNumber, err := getRowNumberByName(name)
	if err != nil {
		return &user.DeleteUserResponse{Success: false}, err
	}

	// Delete the row in the Google Sheet
	//deleteRange := fmt.Sprintf("Sheet1!A%d:E%d", rowNumber, rowNumber) // Assuming data is in columns A-E
	deleteRequest := &sheets.Request{
		DeleteDimension: &sheets.DeleteDimensionRequest{
			Range: &sheets.DimensionRange{
				SheetId:    0, // Assuming you're working with the first sheet
				Dimension:  "ROWS",
				StartIndex: int64(rowNumber - 1), // -1 because sheet indexing starts at 0
				EndIndex:   int64(rowNumber),
			},
		},
	}

	batchUpdateRequest := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{deleteRequest},
	}

	_, err = srv.Spreadsheets.BatchUpdate(spreadsheetID, batchUpdateRequest).Context(ctx).Do()
	if err != nil {
		return &user.DeleteUserResponse{Success: false}, err
	}

	return &user.DeleteUserResponse{Success: true}, nil
}

func (s *userServiceServer) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.User, error) {
	newUser := req.GetUser()

	// Prepare the data to be written to Google Sheets
	var rowData []interface{}
	rowData = append(rowData, newUser.Name)
	rowData = append(rowData, newUser.Age)
	rowData = append(rowData, newUser.CommuteMethod)
	rowData = append(rowData, newUser.College)
	rowData = append(rowData, newUser.Hobbies)

	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{rowData},
	}

	// Append the data to Google Sheets
	_, err := srv.Spreadsheets.Values.Append(spreadsheetID, writeRange, valueRange).ValueInputOption("RAW").Do()
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	user.RegisterUserServiceServer(s, &userServiceServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
