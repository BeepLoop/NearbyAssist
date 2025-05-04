package review_service

import (
	"database/sql"
	"errors"
	"nearbyassist/internal/models"
	booking_repo "nearbyassist/internal/repository/booking"
	review_repo "nearbyassist/internal/repository/review"
	service_repo "nearbyassist/internal/repository/service"
	"nearbyassist/internal/request"
	"nearbyassist/internal/response"
	"nearbyassist/internal/service/core"
	"nearbyassist/internal/utils"
	"slices"
)

const (
	ERR_NOT_FOUND = "not found"
)

type Service struct {
	bookingStore booking_repo.BookingRepository
	reviewStore  review_repo.ReviewRepository
	serviceStore service_repo.ServiceRepository
	encrypt      core.Encryption
	jwt          core.Authenticator
}

func NewService(
	bookingStore booking_repo.BookingRepository,
	reviewStore review_repo.ReviewRepository,
	serviceStore service_repo.ServiceRepository,
	encrypt core.Encryption,
	jwt core.Authenticator,
) *Service {
	return &Service{
		bookingStore: bookingStore,
		reviewStore:  reviewStore,
		serviceStore: serviceStore,
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
		RevieweeId: userId,
		BookingId:  req.BookingId,
		Rating:     req.Rating,
		Text:       utils.Must(s.encrypt.EncryptString(req.Text)),
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

func (s *Service) GetReviewOnBooking(userId, bookingId string) (*models.ReviewModel, error) {
	review, err := s.reviewStore.GetUserReviewOnBooking(userId, bookingId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(ERR_NOT_FOUND)
		}

		return nil, err
	}

	review.Text = utils.Must(s.encrypt.DecryptString(review.Text))
	review.Reviewee.Name = utils.Must(s.encrypt.DecryptString(review.Reviewee.Name))

	return review, nil
}

func (s *Service) GetServiceReviews(serviceId string) ([]response.Review, error) {
	reviews, err := s.serviceStore.GetReviews(serviceId)
	if err != nil {
		return nil, err
	}

	response := slices.AppendSeq(
		make([]response.Review, 0),
		utils.Map(reviews, func(review *models.ReviewModel) response.Review {
			return response.Review{
				Id:               review.Id,
				BookingId:        review.BookingId,
				Rating:           review.Rating,
				Text:             utils.Must(s.encrypt.DecryptString(review.Text)),
				CreatedAt:        review.CreatedAt,
				RevieweeName:     utils.Must(s.encrypt.DecryptString(review.Reviewee.Name)),
				RevieweeImageUrl: review.Reviewee.ImageUrl,
			}
		}),
	)

	return response, nil
}
