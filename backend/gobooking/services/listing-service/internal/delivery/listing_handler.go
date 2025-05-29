package delivery

import (
	"context"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/listing-service/internal/service"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/listing-service/listing"
)

type ListingHandler struct {
	listing.UnimplementedListingServiceServer
	Service *service.ListingService
}

func NewListingHandler(s *service.ListingService) *ListingHandler {
	return &ListingHandler{Service: s}
}

func (h *ListingHandler) CreateListing(ctx context.Context, req *listing.CreateListingRequest) (*listing.ListingResponse, error) {
	l, err := h.Service.Create(ctx, req.Title, req.Description, req.City, req.Price, req.OwnerId, req.Category)
	if err != nil {
		return nil, err
	}

	return &listing.ListingResponse{
		Id:          l.ID,
		Title:       l.Title,
		Description: l.Description,
		City:        l.City,
		Price:       l.Price,
		OwnerId:     l.OwnerID,
		Category:    l.Category,
	}, nil
}

func (h *ListingHandler) GetListingById(ctx context.Context, req *listing.GetListingRequest) (*listing.ListingResponse, error) {
	l, err := h.Service.GetByID(ctx, req.ListingId)
	if err != nil {
		return nil, err
	}
	return &listing.ListingResponse{
		Id:          l.ID,
		Title:       l.Title,
		Description: l.Description,
		City:        l.City,
		Price:       l.Price,
		OwnerId:     l.OwnerID,
		Category:    l.Category,
	}, nil
}

func (h *ListingHandler) SearchListings(ctx context.Context, req *listing.SearchListingsRequest) (*listing.SearchListingsResponse, error) {
	results, err := h.Service.Search(ctx, req.City, req.MinPrice, req.MaxPrice, req.Category)
	if err != nil {
		return nil, err
	}

	var listings []*listing.ListingResponse
	for _, l := range results {
		listings = append(listings, &listing.ListingResponse{
			Id:          l.ID,
			Title:       l.Title,
			Description: l.Description,
			City:        l.City,
			Price:       l.Price,
			OwnerId:     l.OwnerID,
			Category:    l.Category,
		})
	}
	return &listing.SearchListingsResponse{Listings: listings}, nil
}
