package main

import (
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/booking-service/booking"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/booking-service/config"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/booking-service/internal/delivery"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/booking-service/internal/repository"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/booking-service/internal/service"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/booking-service/shared"
	"log"
	"net"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	_ "go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

func main() {
	// Load .env (если есть)
	_ = godotenv.Load()

	// Load config
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Config load error: %v", err)
	}

	// Mongo
	mongoClient, err := shared.InitMongo(cfg.MongoURI)
	if err != nil {
		log.Fatalf("Mongo init error: %v", err)
	}
	db := mongoClient.Database(cfg.MongoDB)

	// Redis
	if err := shared.InitRedis(cfg.RedisAddr); err != nil {
		log.Fatalf("Redis error: %v", err)
	}

	// NATS
	nc, err := nats.Connect(cfg.NatsURL)
	if err != nil {
		log.Fatalf("NATS error: %v", err)
	}
	defer nc.Close()

	// Layer init
	repo := repository.NewBookingRepo(db)
	svc := service.NewBookingService(repo, nc)
	handler := delivery.NewBookingHandler(svc)

	// gRPC
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	booking.RegisterBookingServiceServer(grpcServer, handler)

	log.Println("BookingService running on :50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC serve error: %v", err)
	}
}
