package main

import (
	"log"
	"net/http"
	"time"

	authpb "github.com/fr1gn/bookingsystem/backend/gobooking/api-gateway/auth"
	bookingpb "github.com/fr1gn/bookingsystem/backend/gobooking/api-gateway/booking"
	"github.com/fr1gn/bookingsystem/backend/gobooking/api-gateway/handler"
	listingpb "github.com/fr1gn/bookingsystem/backend/gobooking/api-gateway/listing"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	authAddress := "auth-service:50051"
	bookingAddress := "booking-service:50052"
	listingAddress := "listing-service:50053"

	// gRPC clients
	authConn, err := grpc.Dial(authAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to auth service: %v", err)
	}
	defer authConn.Close()

	bookingConn, err := grpc.Dial(bookingAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to booking service: %v", err)
	}
	defer bookingConn.Close()

	listingConn, err := grpc.Dial(listingAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to listing service: %v", err)
	}
	defer listingConn.Close()

	authClient := authpb.NewAuthServiceClient(authConn)
	bookingClient := bookingpb.NewBookingServiceClient(bookingConn)
	listingClient := listingpb.NewListingServiceClient(listingConn)

	// Gin router
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	handler.RegisterAuthRoutes(r, authClient)
	handler.RegisterBookingRoutes(r, bookingClient)
	handler.RegisterListingRoutes(r, listingClient)

	log.Println("API Gateway running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
