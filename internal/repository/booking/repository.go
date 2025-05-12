package booking_repo

import (
	"nearbyassist/internal/models"
	"time"
)

type BookingRepository interface {
	Create(data *models.BookingModel) (string, error)
	FindById(id string) (*models.BookingModel, error)
	GetAll(limit, offset int) ([]*models.BookingModel, error)
	HasOngoingBookingForService(data *models.BookingModel) (bool, error)

	GetConfirmedBookingsOfVendor(vendorId string) ([]*models.BookingModel, error)

	GetBookingSent(id string) ([]*models.BookingModel, error)
	GetBookingReceived(id string) ([]*models.BookingModel, error)

	GetRecent(id string) ([]*models.BookingModel, error)
	GetConfirmed(id, filter string) ([]*models.BookingModel, error)

	GetHistory(id, filter string) ([]*models.BookingModel, error)
	GetReviewableBookings(userId string) ([]*models.BookingModel, error)

	Cancel(bookingId, cancelledBy, reason string) error
	Accept(bookingId string, scheduleStart, scheduleEnd time.Time) error
	Reject(bookingId, reason string) error
	MarkComplete(bookingId string) error
	Reschedule(bookingId string, scheduleStart, scheduleEnd time.Time) error

	IsReviewed(bookingId string) (bool, error)
}
