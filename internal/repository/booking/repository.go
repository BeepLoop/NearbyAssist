package booking_repo

import (
	"nearbyassist/internal/models"
)

type BookingRepository interface {
	Create(data *models.BookingModel) (string, error)
	FindById(id string) (*models.BookingModel, error)
	GetAll() ([]*models.BookingModel, error)
	HasOngoingBookingForService(data *models.BookingModel) (bool, error)

	GetConfirmedBookingsOfVendor(vendorId string) ([]*models.BookingModel, error)

	GetBookingSent(id string) ([]*models.BookingModel, error)
	GetBookingReceived(id string) ([]*models.BookingModel, error)

	GetRecent(id string) ([]*models.BookingModel, error)
	GetConfirmed(id, filter string) ([]*models.BookingModel, error)

	GetHistory(id, filter string) ([]*models.BookingModel, error)
	GetReviewableBookings(userId string) ([]*models.BookingModel, error)

	Cancel(bookingId, reason string) error
	Accept(bookingId, schedule string) error
	Reject(bookingId, reason string) error
	MarkComplete(bookingId string) error

	IsReviewed(bookingId string) (bool, error)
}
