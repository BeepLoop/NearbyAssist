package middleware

import (
	admin_repo "nearbyassist/internal/repository/admin"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Redirect the user away from reset if not required to change pass
func CheckIfShouldChangePass(repo admin_repo.AdminRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			activeSession, err := utils.GetAdminFromSession(c)
			if err != nil {
				if err := utils.SetFlashMessage(c, "error", "invalid session"); err != nil {
					return c.Redirect(http.StatusSeeOther, "/auth/login?error=invalid_session_error")
				}

				return c.Redirect(http.StatusSeeOther, "/auth/login")
			}

			admin, err := repo.FindById(activeSession.Id)
			if err != nil {
				if err := utils.SetFlashMessage(c, "error", "Unknown user"); err != nil {
					return c.Redirect(http.StatusSeeOther, "/auth/login?error=unknown_users_session")
				}

				return c.Redirect(http.StatusSeeOther, "/auth/login")
			}

			if !admin.MustChangePassword {
				return c.Redirect(http.StatusSeeOther, "/")
			}

			return next(c)
		}
	}
}
