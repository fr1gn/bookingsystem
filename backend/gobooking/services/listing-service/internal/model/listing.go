package model

type Listing struct {
	ID          string  `bson:"_id,omitempty"`
	Title       string  `bson:"title"`
	Description string  `bson:"description"`
	City        string  `bson:"city"`
	Price       float64 `bson:"price"`
	OwnerID     string  `bson:"owner_id"`
	Category    string  `bson:"category"` // apartment, studio, etc.
}
