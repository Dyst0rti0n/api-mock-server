# Advanced Mock Server

This is an advanced API mock server built with Go. It supports CRUD operations and includes features like logging, CORS, rate limiting, input sanitization, JWT-based authentication, role-based access control (RBAC), and Swagger documentation.

## Prerequisites

- Go 1.22.5
- Docker (optional, for Docker support)

## Setup

1. Clone the repository:
    ```sh
    git clone https://github.com/your-repo/advanced-mock-server.git
    cd advanced-mock-server
    ```

2. Install dependencies:
    ```sh
    go mod download
    ```

3. Run the server:
    ```sh
    go run cmd/server/main.go
    ```

## Running with Docker

1. Build the Docker image:
    ```sh
    docker build -t advanced-mock-server .
    ```

2. Run the Docker container:
    ```sh
    docker run -p 8080:8080 advanced-mock-server
    ```

## API Endpoints

Swagger documentation is available at: `http://localhost:8080/swagger/index.html`

## Authentication

Use the `/api/v1/auth` endpoint to obtain a JWT token. Use the token to access protected endpoints by adding it to the `Authorization` header as `Bearer <token>`.

## Testing

Run the tests:
    ```sh
    go test -v ./...
    ```

## Configuration

Configuration is managed through JSON files in the `config` directory. Use environment variables to switch between configurations:
- `config.development.json`
- `config.production.json`
