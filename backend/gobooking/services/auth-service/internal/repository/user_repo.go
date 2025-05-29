package repository

import (
	"auth-service/internal/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	return err
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.Collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	return &user, err
}
