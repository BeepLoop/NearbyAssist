package review_service

import (
	"nearbyassist/internal/models"
	review_repo "nearbyassist/internal/repository/review"
	"nearbyassist/internal/request"
	"nearbyassist/internal/service/auth"
	"nearbyassist/internal/utils"
)

type Service struct {
	store   review_repo.ReviewRepository
	encrypt auth.Encryption
	jwt     auth.Authenticator
}

func NewService(store review_repo.ReviewRepository, encrypt auth.Encryption, jwt auth.Authenticator) *Service {
	return &Service{store: store, encrypt: encrypt, jwt: jwt}
}

func (s *Service) CreateReview(bearerToken string, req *request.NewReviewPayload) (string, error) {
	review, err := s.store.FindById(req.TransactionId)
	if err != nil {
		return "", err
	}

	if err := s.store.IsServiceReviewable(review.ServiceId); err != nil {
		return "", err
	}

	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return "", err
	}

	transaction, err := s.store.GetTransactionById(req.TransactionId)
	if err != nil {
		return "", err
	}

	if transaction.ClientId != userId {
		return "", err
	}

	newReview := new(models.ReviewModel)
	newReview.TransactionId = req.TransactionId
	newReview.ServiceId = req.ServiceId
	newReview.Rating = req.Rating

	reviewId, err := s.store.Create(newReview)
	if err != nil {
		return "", err
	}

	return reviewId, nil
}

func (s *Service) GetReview(reviewId string) (*models.ReviewModel, error) {
	review, err := s.store.FindById(reviewId)
	if err != nil {
		return nil, err
	}

	return review, nil
}

func (s *Service) GetServiceReviews(serviceId string) ([]*models.ReviewModel, error) {
	reviews, err := s.store.GetReviewsByService(serviceId)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}
