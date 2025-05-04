package review_repo

import "nearbyassist/internal/models"

type ReviewRepository interface {
	Create(data *models.ReviewModel) (string, error)
	FindById(id string) (*models.ReviewModel, error)
	GetUserReviewOnBooking(userId, bookingId string) (*models.ReviewModel, error)
}
