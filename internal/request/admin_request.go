package request

type AdminRole string

const (
	ADMIN_ROLE_ADMIN AdminRole = "admin"
	ADMIN_ROLE_STAFF AdminRole = "staff"
)

type NewAdminPayload struct {
	Username string    `json:"username" validate:"required"`
	Password string    `json:"password" validate:"required"`
	Role     AdminRole `json:"role" validate:"required"`
}
