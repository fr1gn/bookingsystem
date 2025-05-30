package repository

import (
	"context"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/auth-service/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type UserRepo struct {
	Collection *mongo.Collection
}

func NewUserRepo(db *mongo.Database) *UserRepo {
	return &UserRepo{
		Collection: db.Collection("users"),
	}
}

func (r *UserRepo) CreateUser(ctx context.Context, user *model.User) error {
	_, err := r.Collection.InsertOne(ctx, user)
	if err != nil {
		log.Println("InsertOne error:", err)
	}
	return err
}

func (r *UserRepo) VerifyEmail(ctx context.Context, email string) error {
	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"email_verified": true}}
	_, err := r.Collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.Collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		log.Println("FindOne error:", err)
	}
	return &user, err
}
