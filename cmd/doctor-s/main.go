package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alikhan-s/doctor-service/internal/repository"
	transport "github.com/alikhan-s/doctor-service/internal/transport/grpc"
	"github.com/alikhan-s/doctor-service/internal/usecase"
	pb "github.com/alikhan-s/doctor-service/proto"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database("doctor_db")

	repo := repository.NewDoctorMongoRepo(db)
	usecaseLayer := usecase.NewDoctorUseCase(repo)

	grpcServer := grpc.NewServer()
	handler := transport.NewDoctorHandler(usecaseLayer)
	pb.RegisterDoctorServiceServer(grpcServer, handler)

	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	reflection.Register(grpcServer)

	go func() {
		log.Println("Doctor gRPC Service is running on port 8081")
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down Doctor gRPC server...")

	grpcServer.GracefulStop()
	log.Println("Server exiting")
}
