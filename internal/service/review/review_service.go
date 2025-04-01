package review_service

import (
	"errors"
	"nearbyassist/internal/models"
	review_repo "nearbyassist/internal/repository/review"
	"nearbyassist/internal/request"
	"nearbyassist/internal/service/core"
	"nearbyassist/internal/utils"
)

type Service struct {
	store   review_repo.ReviewRepository
	encrypt core.Encryption
	jwt     core.Authenticator
}

func NewService(store review_repo.ReviewRepository, encrypt core.Encryption, jwt core.Authenticator) *Service {
	return &Service{store: store, encrypt: encrypt, jwt: jwt}
}

func (s *Service) CreateReview(bearerToken string, req *request.NewReviewPayload) (string, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return "", err
	}

	if isReviewed, err := s.store.IsReviewed(req.ServiceId); err != nil {
		return "", err
	} else {
		if isReviewed {
			return "", errors.New("service already reviewed")
		}
	}

	transaction, err := s.store.GetTransactionById(req.TransactionId)
	if err != nil {
		return "", err
	}

	if transaction.ClientId != userId {
		return "", err
	}

	newReview := &models.ReviewModel{
		TransactionId: req.TransactionId,
		ServiceId:     req.ServiceId,
		Rating:        req.Rating,
		Text:          utils.Must(s.encrypt.EncryptString(req.Text)),
	}

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
