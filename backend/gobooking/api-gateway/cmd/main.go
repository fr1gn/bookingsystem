package main

import (
	"github.com/fr1gn/bookingsystem/backend/gobooking/api-gateway/config"
	"github.com/fr1gn/bookingsystem/backend/gobooking/api-gateway/handler"
	"log"
	"net/http"

	authpb "github.com/fr1gn/bookingsystem/backend/gobooking/services/auth-service/auth"
	bookingpb "github.com/fr1gn/bookingsystem/backend/gobooking/services/booking-service/booking"
	listingpb "github.com/fr1gn/bookingsystem/backend/gobooking/services/listing-service/listing"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.Load("../config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// gRPC clients
	authConn, _ := grpc.Dial(cfg.AuthAddress, grpc.WithInsecure())
	bookingConn, _ := grpc.Dial(cfg.BookingAddress, grpc.WithInsecure())
	listingConn, _ := grpc.Dial(cfg.ListingAddress, grpc.WithInsecure())

	authClient := authpb.NewAuthServiceClient(authConn)
	bookingClient := bookingpb.NewBookingServiceClient(bookingConn)
	listingClient := listingpb.NewListingServiceClient(listingConn)

	// Gin
	r := gin.Default()
	handler.RegisterAuthRoutes(r, authClient)
	handler.RegisterBookingRoutes(r, bookingClient)
	handler.RegisterListingRoutes(r, listingClient)

	log.Println("API Gateway running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
