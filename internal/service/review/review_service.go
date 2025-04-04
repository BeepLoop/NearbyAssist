package review_service

import (
	"errors"
	"nearbyassist/internal/models"
	booking_repo "nearbyassist/internal/repository/booking"
	review_repo "nearbyassist/internal/repository/review"
	"nearbyassist/internal/request"
	"nearbyassist/internal/service/core"
	"nearbyassist/internal/utils"
)

type Service struct {
	bookingStore booking_repo.BookingRepository
	reviewStore  review_repo.ReviewRepository
	encrypt      core.Encryption
	jwt          core.Authenticator
}

func NewService(
	bookingStore booking_repo.BookingRepository,
	reviewStore review_repo.ReviewRepository,
	encrypt core.Encryption,
	jwt core.Authenticator,
) *Service {
	return &Service{
		bookingStore: bookingStore,
		reviewStore:  reviewStore,
		encrypt:      encrypt,
		jwt:          jwt,
	}
}

func (s *Service) CreateReview(bearerToken string, req *request.NewReviewPayload) (string, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return "", err
	}

	booking, err := s.bookingStore.FindById(req.BookingId)
	if err != nil {
		return "", err
	}

	if booking.ClientId != userId {
		return "", err
	}

	if booking.Status != models.BOOKING_STATUS_DONE {
		return "", errors.New("forbidden")
	}

	if isReviewed, err := s.bookingStore.IsReviewed(req.BookingId); err != nil {
		return "", err
	} else {
		if isReviewed {
			return "", errors.New("forbidden")
		}
	}

	newReview := &models.ReviewModel{
		BookingId: req.BookingId,
		Rating:    req.Rating,
		Text:      utils.Must(s.encrypt.EncryptString(req.Text)),
	}

	reviewId, err := s.reviewStore.Create(newReview)
	if err != nil {
		return "", err
	}

	return reviewId, nil
}

func (s *Service) GetReview(reviewId string) (*models.ReviewModel, error) {
	review, err := s.reviewStore.FindById(reviewId)
	if err != nil {
		return nil, err
	}

	return review, nil
}
