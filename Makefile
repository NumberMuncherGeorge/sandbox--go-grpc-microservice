.PHONY: proto build test run clean

# Generate protobuf files
proto:
	protoc -I api/proto -I third_party/googleapis \
		--go_out=api/gen/go/rating/v1 --go_opt=paths=source_relative \
		--go-grpc_out=api/gen/go/rating/v1 --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=api/gen/go/rating/v1 --grpc-gateway_opt=paths=source_relative \
		api/proto/rating.proto

# Build the application
build:
	go build -o bin/rating-service ./cmd/server

# Run tests
test:
	go test -v ./...

# Run the application locally
run:
	go run ./cmd/server

# Clean build artifacts
clean:
	rm -rf bin/

# Tidy modules
tidy:
	go mod tidy

# Lint the code
lint:
	golangci-lint run ./...

# Docker build
docker-build:
	docker build -t rating-service .

# Docker compose up
docker-up:
	docker-compose up -d

# Docker compose down
docker-down:
	docker-compose down
