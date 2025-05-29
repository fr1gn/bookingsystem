package main

import (
	_ "context"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/auth-service/auth"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/auth-service/internal/delivery"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/auth-service/internal/repository"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/auth-service/internal/service"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/auth-service/shared"
	"log"
	"net"

	_ "go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

func main() {

	mongoURI := "mongodb://mongo:27017"
	mongoDB := "gobooking"
	redisAddr := "redis:6379"

	// Connect to Mongo
	mongoClient, err := shared.InitMongo(mongoURI)
	if err != nil {
		log.Fatal("Mongo connect error:", err)
	}
	db := mongoClient.Database(mongoDB)

	// Init Redis
	if err := shared.InitRedis(redisAddr); err != nil {
		log.Fatal("Redis error:", err)
	}

	// Init repository, service, handler
	userRepo := repository.NewUserRepo(db)
	authService := service.NewAuthService(userRepo)
	authHandler := delivery.NewAuthHandler(authService)

	// Start gRPC server
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	auth.RegisterAuthServiceServer(grpcServer, authHandler)

	log.Println("AuthService running on :50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
