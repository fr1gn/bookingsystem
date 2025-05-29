package repository

import (
	"context"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/booking-service/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingRepo struct {
	Collection *mongo.Collection
}

func NewBookingRepo(db *mongo.Database) *BookingRepo {
	return &BookingRepo{
		Collection: db.Collection("bookings"),
	}
}

func (r *BookingRepo) Create(ctx context.Context, b *model.Booking) error {
	_, err := r.Collection.InsertOne(ctx, b)
	return err
}

func (r *BookingRepo) Cancel(ctx context.Context, bookingID string) error {
	_, err := r.Collection.UpdateOne(ctx, bson.M{"_id": bookingID}, bson.M{"$set": bson.M{"status": "cancelled"}})
	return err
}

func (r *BookingRepo) GetByID(ctx context.Context, id string) (*model.Booking, error) {
	var booking model.Booking
	err := r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&booking)
	return &booking, err
}

func (r *BookingRepo) ListByUser(ctx context.Context, userID string) ([]*model.Booking, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []*model.Booking
	for cursor.Next(ctx) {
		var b model.Booking
		if err := cursor.Decode(&b); err == nil {
			bookings = append(bookings, &b)
		}
	}
	return bookings, nil
}
