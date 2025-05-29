package model

type Booking struct {
	ID           string `bson:"_id,omitempty"`
	UserID       string `bson:"user_id"`
	ListingID    string `bson:"listing_id"`
	StartDate    string `bson:"start_date"`
	EndDate      string `bson:"end_date"`
	DurationType string `bson:"duration_type"` // "hour", "day", "month"
	Status       string `bson:"status"`        // "active", "cancelled"
}
