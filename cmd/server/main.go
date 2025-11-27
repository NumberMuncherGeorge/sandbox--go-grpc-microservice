// Package main implements the gRPC server with HTTP transcoding for the rating service.
package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	pb "github.com/NumberMuncherGeorge/sandbox--go-grpc-microservice/api/gen/go/rating/v1"
	"github.com/NumberMuncherGeorge/sandbox--go-grpc-microservice/internal/repository"
	"github.com/NumberMuncherGeorge/sandbox--go-grpc-microservice/internal/service"
	"github.com/NumberMuncherGeorge/sandbox--go-grpc-microservice/pkg/config"
)

func main() {
	cfg := config.Load()

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	// Ping MongoDB to verify connection
	if err := mongoClient.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	log.Println("Connected to MongoDB")

	db := mongoClient.Database(cfg.MongoDB)

	// Create repositories
	resourceRepo := repository.NewResourceRepository(db)
	ratingRepo := repository.NewRatingRepository(db)

	// Create service
	ratingService := service.NewRatingService(resourceRepo, ratingRepo)

	// Start gRPC server
	go func() {
		if err := startGRPCServer(cfg.GRPCPort, ratingService); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	// Start HTTP gateway
	if err := startHTTPGateway(cfg.HTTPPort, cfg.GRPCPort); err != nil {
		log.Fatalf("Failed to start HTTP gateway: %v", err)
	}
}

func startGRPCServer(port string, ratingService *service.RatingService) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRatingServiceServer(grpcServer, ratingService)

	// Enable reflection for gRPC CLI tools
	reflection.Register(grpcServer)

	log.Printf("gRPC server listening on port %s", port)
	return grpcServer.Serve(lis)
}

func startHTTPGateway(httpPort, grpcPort string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := pb.RegisterRatingServiceHandlerFromEndpoint(ctx, mux, "localhost:"+grpcPort, opts)
	if err != nil {
		return err
	}

	log.Printf("HTTP gateway listening on port %s", httpPort)
	return http.ListenAndServe(":"+httpPort, mux)
}
