# Service Registry Server

This Go-based Service Registry server allows microservices within an architecture to register themselves and their status to be monitored periodically. The service exposes a REST API for service registration and querying service details.

## Features

- **Service Registration**: Microservices can self-register with their metadata like name, host, port, and service type.
- **Service Discovery**: Clients or other microservices can retrieve a list of all registered services with their current status and details.
- **Health Check Monitoring**: The registry performs periodic health checks on each registered service to update their status.

## Getting Started

### Prerequisites

- Install [Go](https://golang.org/doc/install) (version 1.15 or higher recommended).
- Install protobuf
- Ensure you have generated the necessary Protocol Buffer files for `pb.APIService` and `pb.Registration`.

Run this command to generate the pb.go files:

```bash
    protoc --go_out=. service_registry.proto
```


### Installation

1. Clone the repository to your local machine:

    ```bash
    git clone https://github.com/<YourUsername>/<YourRepositoryName>.git
    cd <YourRepositoryName>
    ```

2. Run the Go server:

    ```bash
    go run main.go
    ```

### Register a Service

Send a `POST` request to `/register` with a JSON body containing the name, host, port, and type of the service.

Example:

```json
{
  "Name": "example-service",
  "Host": "localhost",
  "Port": 8080,
  "Type": "HTTP"
}
```

### Get List of Services

Send a `GET` request to `/services` to obtain an array of JSON objects containing the name, status, host, port, and type of the service. 
