package user_management_service

func (s *Service) BanUser(userId string) error {
	if err := s.userStore.BanUser(userId); err != nil {
		return err
	}

	// TODO: Notify user

	return nil
}

func (s *Service) UnbanUser(userId string) error {
	if err := s.userStore.UnbanUser(userId); err != nil {
		return err
	}

	// TODO: Notify user

	return nil
}
