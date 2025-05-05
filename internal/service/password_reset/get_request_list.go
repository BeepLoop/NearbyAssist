package passwordreset_service

import (
	"nearbyassist/internal/dto"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
	"slices"
)

func (s *Service) GetRequestList() ([]dto.PasswordResetRequest, error) {
	requests, err := s.passwordResetStore.GetAll()
	if err != nil {
		return nil, err
	}

	data := slices.AppendSeq(
		make([]dto.PasswordResetRequest, 0),
		utils.Map(requests, func(req *models.PasswordResetRequestModel) dto.PasswordResetRequest {
			return dto.PasswordResetRequest{
				Id:        req.Id,
				AdminId:   req.AdminId,
				Username:  utils.Must(s.encrypt.DecryptString(req.Username)),
				Email:     utils.Must(s.encrypt.DecryptString(req.Email)),
				CreatedAt: utils.FormatDate(req.CreatedAt),
			}
		}),
	)

	return data, nil
}
