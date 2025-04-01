package review_service

import (
	"errors"
	"nearbyassist/internal/models"
	review_repo "nearbyassist/internal/repository/review"
	transaction_repo "nearbyassist/internal/repository/transaction"
	"nearbyassist/internal/request"
	"nearbyassist/internal/service/core"
	"nearbyassist/internal/utils"
)

type Service struct {
	transactionStore transaction_repo.TransactionRepository
	reviewStore      review_repo.ReviewRepository
	encrypt          core.Encryption
	jwt              core.Authenticator
}

func NewService(
	transactionStore transaction_repo.TransactionRepository,
	reviewStore review_repo.ReviewRepository,
	encrypt core.Encryption,
	jwt core.Authenticator,
) *Service {
	return &Service{
		transactionStore: transactionStore,
		reviewStore:      reviewStore,
		encrypt:          encrypt,
		jwt:              jwt,
	}
}

func (s *Service) CreateReview(bearerToken string, req *request.NewReviewPayload) (string, error) {
	userId, err := utils.GetUserIdFromToken(bearerToken, s.jwt.GetClaims)
	if err != nil {
		return "", err
	}

	if isReviewable, err := s.transactionStore.IsReviewable(req.TransactionId); err != nil {
		return "", err
	} else {
		if !isReviewable {
			return "", errors.New("transaction not reviewable")
		}
	}

	if isReviewed, err := s.reviewStore.IsReviewed(req.ServiceId); err != nil {
		return "", err
	} else {
		if isReviewed {
			return "", errors.New("service already reviewed")
		}
	}

	transaction, err := s.reviewStore.GetTransactionById(req.TransactionId)
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

func (s *Service) GetServiceReviews(serviceId string) ([]*models.ReviewModel, error) {
	reviews, err := s.reviewStore.GetReviewsByService(serviceId)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}
