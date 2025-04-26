package passwordreset_service

func (s *Service) RejectResetPassword(requestId, reason string) error {
	// NOTE: Implement reason
	return s.passwordResetStore.Delete(requestId)
}
