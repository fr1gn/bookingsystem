package handler

import (
	"github.com/fr1gn/bookingsystem/backend/gobooking/api-gateway/booking"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterBookingRoutes(r *gin.Engine, client booking.BookingServiceClient) {
	r.POST("/booking/create", func(c *gin.Context) {
		var req struct {
			UserID       string `json:"user_id"`
			ListingID    string `json:"listing_id"`
			StartDate    string `json:"start_date"`
			EndDate      string `json:"end_date"`
			DurationType string `json:"duration_type"`
		}

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, err := client.CreateBooking(c, &booking.CreateBookingRequest{
			UserId:       req.UserID,
			ListingId:    req.ListingID,
			StartDate:    req.StartDate,
			EndDate:      req.EndDate,
			DurationType: req.DurationType,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
	})
}
