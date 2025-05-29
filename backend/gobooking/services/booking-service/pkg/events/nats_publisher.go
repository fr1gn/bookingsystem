package events

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
)

type BookingCreatedEvent struct {
	BookingID string `json:"booking_id"`
	UserID    string `json:"user_id"`
	ListingID string `json:"listing_id"`
}

func PublishBookingCreated(nc *nats.Conn, event BookingCreatedEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return nc.Publish("booking.created", data)
}
