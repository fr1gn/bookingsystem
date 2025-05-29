package repository

import (
	"context"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/listing-service/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ListingRepo struct {
	Collection *mongo.Collection
}

func NewListingRepo(db *mongo.Database) *ListingRepo {
	return &ListingRepo{
		Collection: db.Collection("listings"),
	}
}

func (r *ListingRepo) Create(ctx context.Context, listing *model.Listing) error {
	_, err := r.Collection.InsertOne(ctx, listing)
	return err
}

func (r *ListingRepo) GetByID(ctx context.Context, id string) (*model.Listing, error) {
	var l model.Listing
	err := r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&l)
	return &l, err
}

func (r *ListingRepo) Search(ctx context.Context, city string, min, max float64, category string) ([]*model.Listing, error) {
	filter := bson.M{}
	if city != "" {
		filter["city"] = city
	}
	if category != "" {
		filter["category"] = category
	}
	filter["price"] = bson.M{"$gte": min, "$lte": max}

	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []*model.Listing
	for cursor.Next(ctx) {
		var l model.Listing
		if err := cursor.Decode(&l); err == nil {
			results = append(results, &l)
		}
	}
	return results, nil
}
