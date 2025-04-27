package admin_service

import (
	"nearbyassist/internal/dto"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
	"slices"
)

func (s *Service) GetAccounts(role string) ([]dto.Admin, error) {
	accounts := make([]*models.AdminModel, 0)

	if role == "" || role == "all" {
		if res, err := s.adminStore.GetAll(); err != nil {
			return nil, err
		} else {
			accounts = res
		}
	} else {
		if res, err := s.adminStore.GetAllWithRole(role); err != nil {
			return nil, err
		} else {
			accounts = res
		}
	}

	data := slices.AppendSeq(
		make([]dto.Admin, 0),
		utils.Map(accounts, func(account *models.AdminModel) dto.Admin {
			return dto.Admin{
				Id:                 account.Id,
				Username:           utils.Must(s.encrypt.DecryptString(account.Username)),
				Email:              utils.Must(s.encrypt.DecryptString(account.Email)),
				Role:               account.Role,
				MustChangePassword: account.MustChangePassword,
				Suspended:          account.Suspended,
				CreatedAt:          utils.FormatDate(account.CreatedAt),
				UpdatedAt:          utils.FormatDate(account.UpdatedAt),
			}
		}),
	)

	return data, nil
}
