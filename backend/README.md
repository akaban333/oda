# Study Platform Backend

This is the backend service for the Study Platform application. It provides a RESTful API for the frontend to interact with.

## Technologies

- **Go** - Programming language
- **Gin** - Web framework
- **MongoDB** - Database
- **JWT** - Authentication
- **MinIO** - Object storage for files
- **WebSocket** - Real-time communication
- **Docker** - Containerization

## Directory Structure

```
backend/
├── .env.example           # Example environment variables
├── cmd/                   # Command-line applications
│   └── api/               # API gateway service
│       └── main.go        # Entry point for the API gateway
├── internal/              # Private application code
│   ├── auth/              # Authentication service
│   ├── room/              # Room management service
│   ├── session/           # Session tracking service
│   ├── content/           # Content management service
│   └── realtime/          # Real-time communication service
├── pkg/                   # Public library code
│   ├── api/               # API utilities
│   ├── auth/              # Authentication utilities
│   ├── database/          # Database connection
│   ├── logger/            # Logging utilities
│   ├── middleware/        # HTTP middleware
│   ├── models/            # Data models
│   ├── storage/           # File storage
│   └── utils/             # Utility functions
└── test/                  # Test files
```

## Setup

### Prerequisites

- Go 1.21 or higher
- MongoDB
- MinIO (for file storage)

### Installation

1. Clone the repository
2. Copy `.env.example` to `.env` and update the values
3. Install dependencies:

```bash
go mod download
```

4. Run the application:

```bash
go run cmd/api/main.go
```

## API Documentation

The API documentation is available at `/api/v1/docs` when the server is running.

## Development

### Running the server

```bash
go run cmd/api/main.go
```

### Building the binary

```bash
go build -o studyplatform cmd/api/main.go
```

### Testing

Run all tests:

```bash
go test ./...
```

## Environment Variables

See the `.env.example` file for a list of required environment variables.

## Docker

Build the Docker image:

```bash
docker build -t studyplatform-backend .
```

Run the Docker container:

```bash
docker run -p 8080:8080 studyplatform-backend
```

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin feature-name`
5. Submit a pull request

## License

This project is licensed under the MIT License. 