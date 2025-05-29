package handler

import (
	"github.com/fr1gn/bookingsystem/backend/gobooking/api-gateway/listing"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func RegisterListingRoutes(r *gin.Engine, client listing.ListingServiceClient) {
	r.POST("/listing/create", func(c *gin.Context) {
		var req struct {
			Title       string  `json:"title"`
			Description string  `json:"description"`
			City        string  `json:"city"`
			Price       float64 `json:"price"`
			OwnerID     string  `json:"owner_id"`
			Category    string  `json:"category"`
		}

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, err := client.CreateListing(c, &listing.CreateListingRequest{
			Title:       req.Title,
			Description: req.Description,
			City:        req.City,
			Price:       req.Price,
			OwnerId:     req.OwnerID,
			Category:    req.Category,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
	})

	r.GET("/listing/search", func(c *gin.Context) {
		city := c.Query("city")
		category := c.Query("category")
		min := parseFloatOrZero(c.Query("min_price"))
		max := parseFloatOrZero(c.Query("max_price"))

		res, err := client.SearchListings(c, &listing.SearchListingsRequest{
			City:     city,
			Category: category,
			MinPrice: min,
			MaxPrice: max,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res.Listings)
	})
}

func parseFloatOrZero(str string) float64 {
	f, _ := strconv.ParseFloat(str, 64)
	return f
}
