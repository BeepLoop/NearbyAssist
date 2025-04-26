package passwordreset_service

const (
	ERR_INVALID_CREDENTIALS     = "invalid credentials"
	ERR_WEAK_PASSWORD           = "password too weak"
	ERR_SELF_RESETTING_PASSWORD = "not allowed to reset own password"
	ERR_MISMATCHING_PASSWORD    = "confirmation password mismatch"
	ERR_INVALID_USERNAME        = "invalid username"
)
