package delivery

import (
	"context"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/booking-service/booking"
	"github.com/fr1gn/bookingsystem/backend/gobooking/services/booking-service/internal/service"
)

type BookingHandler struct {
	booking.UnimplementedBookingServiceServer
	Service *service.BookingService
}

func NewBookingHandler(s *service.BookingService) *BookingHandler {
	return &BookingHandler{Service: s}
}

func (h *BookingHandler) CreateBooking(ctx context.Context, req *booking.CreateBookingRequest) (*booking.BookingResponse, error) {
	b, err := h.Service.CreateBooking(ctx, req.UserId, req.ListingId, req.StartDate, req.EndDate, req.DurationType)
	if err != nil {
		return nil, err
	}

	return &booking.BookingResponse{
		BookingId:    b.ID,
		UserId:       b.UserID,
		ListingId:    b.ListingID,
		StartDate:    b.StartDate,
		EndDate:      b.EndDate,
		DurationType: b.DurationType,
		Status:       b.Status,
	}, nil
}

func (h *BookingHandler) CancelBooking(ctx context.Context, req *booking.CancelBookingRequest) (*booking.BookingResponse, error) {
	err := h.Service.CancelBooking(ctx, req.BookingId)
	if err != nil {
		return nil, err
	}
	return &booking.BookingResponse{BookingId: req.BookingId, Status: "cancelled"}, nil
}

func (h *BookingHandler) GetBookingById(ctx context.Context, req *booking.GetBookingRequest) (*booking.BookingResponse, error) {
	b, err := h.Service.GetBookingByID(ctx, req.BookingId)
	if err != nil {
		return nil, err
	}

	return &booking.BookingResponse{
		BookingId:    b.ID,
		UserId:       b.UserID,
		ListingId:    b.ListingID,
		StartDate:    b.StartDate,
		EndDate:      b.EndDate,
		DurationType: b.DurationType,
		Status:       b.Status,
	}, nil
}

func (h *BookingHandler) ListBookingsForUser(ctx context.Context, req *booking.ListBookingsRequest) (*booking.ListBookingsResponse, error) {
	bookings, err := h.Service.ListBookingsForUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	var res []*booking.BookingResponse
	for _, b := range bookings {
		res = append(res, &booking.BookingResponse{
			BookingId:    b.ID,
			UserId:       b.UserID,
			ListingId:    b.ListingID,
			StartDate:    b.StartDate,
			EndDate:      b.EndDate,
			DurationType: b.DurationType,
			Status:       b.Status,
		})
	}

	return &booking.ListBookingsResponse{Bookings: res}, nil
}
