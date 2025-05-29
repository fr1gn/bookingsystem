package service

import (
	"context"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/listing-service/internal/model"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/listing-service/internal/repository"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/listing-service/pkg/cache"

	"github.com/google/uuid"
)

type ListingService struct {
	Repo  *repository.ListingRepo
	Cache *cache.Cache
}

func NewListingService(repo *repository.ListingRepo, cache *cache.Cache) *ListingService {
	return &ListingService{
		Repo:  repo,
		Cache: cache,
	}
}

func (s *ListingService) Create(ctx context.Context, title, desc, city string, price float64, ownerID, category string) (*model.Listing, error) {
	l := &model.Listing{
		ID:          uuid.NewString(),
		Title:       title,
		Description: desc,
		City:        city,
		Price:       price,
		OwnerID:     ownerID,
		Category:    category,
	}
	err := s.Repo.Create(ctx, l)
	return l, err
}

func (s *ListingService) GetByID(ctx context.Context, id string) (*model.Listing, error) {
	return s.Repo.GetByID(ctx, id)
}

func (s *ListingService) Search(ctx context.Context, city string, min, max float64, category string) ([]*model.Listing, error) {
	return s.Repo.Search(ctx, city, min, max, category)
}
