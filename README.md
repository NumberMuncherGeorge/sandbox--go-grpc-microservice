# Go gRPC Rating Microservice

A Go microservice that provides a gRPC API with HTTP transcoding for managing resources and their ratings. Resources can be created and stored in MongoDB, and users can submit ratings for these resources based on predefined categories.

## Features

- **gRPC API** - High-performance gRPC service for all operations
- **HTTP Transcoding** - RESTful HTTP endpoints via gRPC-Gateway
- **MongoDB Storage** - Persistent storage using MongoDB
- **Resource Management** - Create, retrieve, and list resources
- **Rating System** - Submit and retrieve ratings by category
- **Rating Summaries** - Aggregated statistics per category

## API Endpoints

### gRPC Methods

| Method | Description |
|--------|-------------|
| `CreateResource` | Create a new resource that can be rated |
| `GetResource` | Retrieve a resource by ID |
| `ListResources` | List all resources with pagination |
| `SubmitRating` | Submit a rating for a resource in a category |
| `GetResourceRatings` | Get all ratings and summaries for a resource |

### HTTP Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/v1/resources` | Create a new resource |
| `GET` | `/v1/resources/{id}` | Get a resource by ID |
| `GET` | `/v1/resources` | List all resources |
| `POST` | `/v1/resources/{resource_id}/ratings` | Submit a rating |
| `GET` | `/v1/resources/{resource_id}/ratings` | Get ratings for a resource |

## Project Structure

```
.
├── api/
│   ├── gen/go/rating/v1/     # Generated protobuf Go code
│   └── proto/                 # Protocol buffer definitions
├── cmd/
│   └── server/                # Main server application
├── internal/
│   ├── model/                 # Data models
│   ├── repository/            # MongoDB repositories
│   └── service/               # gRPC service implementations
├── pkg/
│   └── config/                # Configuration management
├── third_party/
│   └── googleapis/            # Google API proto dependencies
├── docker-compose.yml         # Docker Compose for local development
├── Dockerfile                 # Multi-stage Docker build
└── Makefile                   # Build and development commands
```

## Prerequisites

- Go 1.21 or later
- Protocol Buffers compiler (protoc)
- MongoDB (or use Docker Compose)
- gRPC and gRPC-Gateway plugins for protoc

## Installation

### Install Go dependencies

```bash
go mod download
```

### Install protoc plugins

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
```

### Generate protobuf code

```bash
make proto
```

## Running the Service

### Using Docker Compose (recommended)

```bash
docker-compose up -d
```

This starts MongoDB and the rating service. The service will be available at:
- gRPC: `localhost:50051`
- HTTP: `localhost:8080`

### Running locally

1. Start MongoDB:
```bash
docker run -d -p 27017:27017 mongo:7
```

2. Run the service:
```bash
make run
```

Or directly:
```bash
go run ./cmd/server
```

## Configuration

The service can be configured via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `GRPC_PORT` | `50051` | gRPC server port |
| `HTTP_PORT` | `8080` | HTTP gateway port |
| `MONGO_URI` | `mongodb://localhost:27017` | MongoDB connection URI |
| `MONGO_DB` | `rating_service` | MongoDB database name |

## Usage Examples

### Create a Resource (HTTP)

```bash
curl -X POST http://localhost:8080/v1/resources \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Awesome Product",
    "description": "A really great product",
    "categories": ["quality", "value", "design"]
  }'
```

### Get a Resource (HTTP)

```bash
curl http://localhost:8080/v1/resources/{resource_id}
```

### List Resources (HTTP)

```bash
curl http://localhost:8080/v1/resources
```

### Submit a Rating (HTTP)

```bash
curl -X POST http://localhost:8080/v1/resources/{resource_id}/ratings \
  -H "Content-Type: application/json" \
  -d '{
    "category": "quality",
    "score": 5,
    "comment": "Excellent quality!",
    "rater_id": "user123"
  }'
```

### Get Ratings (HTTP)

```bash
curl http://localhost:8080/v1/resources/{resource_id}/ratings
```

### Using gRPC (with grpcurl)

```bash
# List resources
grpcurl -plaintext localhost:50051 rating.v1.RatingService/ListResources

# Create a resource
grpcurl -plaintext -d '{"name": "Test", "description": "Test product", "categories": ["quality"]}' \
  localhost:50051 rating.v1.RatingService/CreateResource
```

## Development

### Build

```bash
make build
```

### Run tests

```bash
make test
```

### Lint (requires golangci-lint)

```bash
make lint
```

### Clean

```bash
make clean
```

## Testing

Tests require a running MongoDB instance. If MongoDB is not available, tests will be skipped.

```bash
# Start MongoDB for testing
docker run -d -p 27017:27017 mongo:7

# Run tests
go test -v ./...
```

## Architecture

The service follows a clean architecture pattern:

1. **API Layer** (`api/proto/`) - Protocol buffer definitions define the service contract
2. **Service Layer** (`internal/service/`) - gRPC handlers implement business logic
3. **Repository Layer** (`internal/repository/`) - MongoDB operations for data persistence
4. **Model Layer** (`internal/model/`) - Internal data structures

The gRPC-Gateway automatically translates HTTP requests to gRPC calls, providing a RESTful API without additional implementation effort.
