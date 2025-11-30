package main

import (
	"fmt"
	"log"
	"net"
	"os"

	airplanev1 "github.com/yourusername/airplane-rating/internal/pb/airplane/v1"
	"github.com/yourusername/airplane-rating/internal/service"
	"google.golang.org/grpc"
)

func main() {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "8080"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	// Initialize services
	airplaneService := service.NewAirplaneService()

	// TODO: Register services here
	airplanev1.RegisterAirplaneServiceServer(grpcServer, airplaneService)
	// categoryv1.RegisterCategoryServiceServer(grpcServer, categoryService)
	// ratingv1.RegisterRatingServiceServer(grpcServer, ratingService)

	log.Printf("Starting gRPC server on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
