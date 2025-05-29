package main

import (
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/listing-service/config"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/listing-service/internal/delivery"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/listing-service/internal/repository"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/listing-service/internal/service"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/listing-service/listing"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/listing-service/pkg/cache"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/listing-service/shared"
	"log"
	"net"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	// Load .env if present
	_ = godotenv.Load()

	// Load config
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Mongo
	mongoClient, err := shared.InitMongo(cfg.MongoURI)
	if err != nil {
		log.Fatalf("failed to connect to mongo: %v", err)
	}
	db := mongoClient.Database(cfg.MongoDB)

	// Redis
	redisClient := shared.InitRedis(cfg.RedisAddr)
	listingCache := cache.NewCache(redisClient)

	// Init layers
	repo := repository.NewListingRepo(db)
	svc := service.NewListingService(repo, listingCache)
	handler := delivery.NewListingHandler(svc)

	// gRPC
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	listing.RegisterListingServiceServer(grpcServer, handler)

	log.Println("ListingService running on :50053")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
