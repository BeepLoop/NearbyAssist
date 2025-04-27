package admin_service

import "errors"

func (s *Service) Suspend(handlerAdminId, targetAccountId string) error {
	admin, err := s.adminStore.FindById(handlerAdminId)
	if err != nil {
		return err
	}
	if admin.Role != "admin" {
		return errors.New(ERR_UNAUTHORIZED)
	}

	target, err := s.adminStore.FindById(targetAccountId)
	if err != nil {
		return err
	}
	if target.Role == "admin" {
		return errors.New(ERR_UNAUTHORIZED)
	}

	if err := s.adminStore.Suspend(targetAccountId); err != nil {
		return err
	}

	return nil
}

func (s *Service) Unsuspend(handlerAdminId, targetAccountId string) error {
	admin, err := s.adminStore.FindById(handlerAdminId)
	if err != nil {
		return err
	}
	if admin.Role != "admin" {
		return errors.New(ERR_UNAUTHORIZED)
	}

	target, err := s.adminStore.FindById(targetAccountId)
	if err != nil {
		return err
	}
	if target.Role == "admin" {
		return errors.New(ERR_UNAUTHORIZED)
	}

	if err := s.adminStore.Unsuspend(targetAccountId); err != nil {
		return err
	}

	return nil
}
