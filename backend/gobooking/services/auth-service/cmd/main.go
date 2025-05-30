package main

import (
	"log"
	"net"
	"os"

	"github.com/fr1gn/bookingsystem/backend/gobooking/services/auth-service/auth"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/auth-service/internal/delivery"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/auth-service/internal/repository"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/auth-service/internal/service"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/auth-service/shared"

	"google.golang.org/grpc"
)

func main() {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	mongoDB := os.Getenv("MONGO_DB")
	if mongoDB == "" {
		mongoDB = "gobooking"
	}
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	mongoClient, err := shared.InitMongo(mongoURI)
	if err != nil {
		log.Fatalf("Mongo connect error: %v", err)
	}
	db := mongoClient.Database(mongoDB)
	log.Println("Connected to MongoDB:", mongoURI)

	if err := shared.InitRedis(redisAddr); err != nil {
		log.Fatalf("Redis connect error: %v", err)
	}
	log.Println("Connected to Redis:", redisAddr)

	userRepo := repository.NewUserRepo(db)
	authService := service.NewAuthService(userRepo)
	authHandler := delivery.NewAuthHandler(authService)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	grpcServer := grpc.NewServer()
	auth.RegisterAuthServiceServer(grpcServer, authHandler)

	log.Println("AuthService is running on :50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
