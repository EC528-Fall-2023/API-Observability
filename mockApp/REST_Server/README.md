# REST-API_GO

# Google Sheets User Management via REST API

This application manages user data by interacting with Google Sheets via the Google Sheets API.

## Overview

The primary functionality includes:
1. Retrieving all users from a specific Google Sheet.
2. Fetching details of a single user based on their name.
3. Deleting a user entry based on their name.
4. Updating user information.
5. Adding a new user.

The user data model includes:
- Name
- Age
- Commute Method
- College
- Hobbies

The application uses the `gin-gonic/gin` package to provide an HTTP server and `google.golang.org/api/sheets/v4` for Google Sheets API interactions.

## Endpoints

1. **Get All Users**: This endpoint fetches all users' details stored in the Google Sheet.
2. **Get User by Name**: Retrieve specific user details based on a given name.
3. **Delete User by Name**: Delete a user entry from the Google Sheet based on the name.
4. **Update User by Name**: Update the details of a user specified by their name.
5. **Add User**: Add a new user to the Google Sheet.

## Setup & Configuration

1. **Credentials**: The application expects a `credentials.json` file to authenticate with Google Sheets. This file contains the necessary credentials to access the Google Sheets API.
2. **Token**: An additional `token.json` is required for OAuth2 authentication. This token is generated upon the first run and used in subsequent API calls.

If you change the scopes in the configuration, ensure you delete the previously saved `token.json`.

3. **Google Sheets Configuration**: The application is hardcoded to interact with a specific Google Sheet identified by `10-CfbfktbeTSMV3tgnIKwaBquzw-RmjS13Tut9A32_s` and specifically the `Sheet1` within that spreadsheet. This is serving as a database to store user information. 

## Code Structure

- The application primarily uses the Gin framework to define and handle HTTP endpoints.
- Google Sheets operations are abstracted into separate functions to create, read, update, and delete user entries.
- Authentication and token management for Google Sheets API are managed through helper functions.

## How to Run

To run the application, ensure you have the credentials and configuration in place, and then simply execute the main Go file.


