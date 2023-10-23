# React Google Sheet App

This project is a React application that interfaces with a server backend to perform CRUD operations on user data. The app allows users to submit, update, delete, search for, and display user information. It also supports switching between REST and gRPC API modes, although the functionality for gRPC is not implemented in the provided code.

## Table of Contents

- [Installation](#installation)
- [Setup](#setup)
- [Running the App](#running-the-app)
- [Usage](#usage)
- [Features](#features)
- [API Endpoints](#api-endpoints)
- [Technologies Used](#technologies-used)
## Installation

Before you can run this application, make sure you have [Node.js and npm](https://nodejs.org/en/download/) installed on your system.

1.  Install the required dependencies:
   ```sh
   npm install
   ```

## Setup

Before starting the application, make sure your backend server is running and accessible. The default URL for the backend is `http://localhost:4000`. You can change this URL in the code if your backend is running on a different address.

## Running the App

Once everything is set up, you can start the application by running:

```sh
npm start
```

This will start the development server and open the application in your default web browser. The application will be available at `http://localhost:3000`.

## Usage

To use the application, fill in the user information in the form provided and use the buttons to perform CRUD operations:

- **Submit:** Add a new user to the database.
- **Update:** Update an existing user's information in the database based on the name field.
- **Delete:** Remove an existing user from the database based on the name field.
- **Search:** Find and display a user's information based on the name field.
- **Show All:** Retrieve and display all users' information from the database.

## Features

- Interactive form for submitting, updating, and deleting user information.
- Table display for showing user data.
- Ability to switch between REST and gRPC API modes.

## API Endpoints

- **POST `/adduser`:** Add a new user.
- **PUT `/updateuser/{name}`:** Update an existing user's information.
- **POST `/deleteuser/{name}`:** Delete an existing user.
- **GET `/getuser/{name}`:** Retrieve a specific user's information.
- **GET `/getuser`:** Retrieve all users' information.

## Technologies Used

- [React](https://reactjs.org/): A JavaScript library for building user interfaces.
- [Semantic UI React](https://react.semantic-ui.com/): A React integration of Semantic UI for styling.
- [Axios](https://axios-http.com/): A promise-based HTTP client for making HTTP requests.
- [Node.js](https://nodejs.org/): A JavaScript runtime built on Chrome's V8 JavaScript engine.

