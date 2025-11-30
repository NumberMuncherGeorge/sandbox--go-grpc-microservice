# Download googleapis for HTTP annotations
FROM alpine:latest as proto-deps
RUN apk add --no-cache git
WORKDIR /proto-deps
RUN git clone --depth 1 --branch master https://github.com/googleapis/googleapis.git

# Build stage
FROM golang:1.25-trixie as builder
WORKDIR /build

# Copy proto dependencies
COPY --from=proto-deps /proto-deps/googleapis /proto-deps/googleapis

# Install protoc and plugins
RUN apt-get update && apt-get install -y protobuf-compiler
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest

# Copy go module files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Generate proto files
RUN mkdir -p internal/pb && \
    protoc \
      --proto_path=proto \
      --proto_path=/proto-deps/googleapis \
      --go_out=internal/pb \
      --go_opt=paths=source_relative \
      --go-grpc_out=internal/pb \
      --go-grpc_opt=paths=source_relative \
      --grpc-gateway_out=internal/pb \
      --grpc-gateway_opt=paths=source_relative \
      --grpc-gateway_opt=generate_unbound_methods=true \
      proto/**/*.proto

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /airplane-rating-server ./cmd/server

# Runtime stage
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app

COPY --from=builder /airplane-rating-server /app/airplane-rating-server

EXPOSE 8080 9090

ENTRYPOINT ["/app/airplane-rating-server"]
