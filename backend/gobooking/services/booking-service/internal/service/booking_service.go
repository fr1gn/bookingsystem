package service

import (
	"context"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/booking-service/internal/model"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/booking-service/internal/repository"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/booking-service/pkg/events"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type BookingService struct {
	Repo *repository.BookingRepo
	NATS *nats.Conn
}

func NewBookingService(repo *repository.BookingRepo, nc *nats.Conn) *BookingService {
	return &BookingService{Repo: repo, NATS: nc}
}

func (s *BookingService) CreateBooking(ctx context.Context, userID, listingID, start, end, dtype string) (*model.Booking, error) {
	booking := &model.Booking{
		ID:           uuid.NewString(),
		UserID:       userID,
		ListingID:    listingID,
		StartDate:    start,
		EndDate:      end,
		DurationType: dtype,
		Status:       "active",
	}

	err := s.Repo.Create(ctx, booking)
	if err != nil {
		return nil, err
	}

	// NATS Event
	go events.PublishBookingCreated(s.NATS, events.BookingCreatedEvent{
		BookingID: booking.ID,
		UserID:    userID,
		ListingID: listingID,
	})

	return booking, nil
}

func (s *BookingService) CancelBooking(ctx context.Context, id string) error {
	return s.Repo.Cancel(ctx, id)
}

func (s *BookingService) GetBookingByID(ctx context.Context, id string) (*model.Booking, error) {
	return s.Repo.GetByID(ctx, id)
}

func (s *BookingService) ListBookingsForUser(ctx context.Context, userID string) ([]*model.Booking, error) {
	return s.Repo.ListByUser(ctx, userID)
}
