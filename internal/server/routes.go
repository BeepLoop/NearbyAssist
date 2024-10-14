package server

import (
	"nearbyassist/views/pages"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) routes() {
	s.Echo.Static("/static", "static")
	s.Echo.Static("/public", "public")

	s.Echo.GET("", func(c echo.Context) error {
		page := pages.Index()
		return s.render(c, http.StatusOK, page)
	})

	s.Echo.GET("/login", func(c echo.Context) error {
		page := pages.Login()
		return s.render(c, http.StatusOK, page)
	})

	s.Echo.POST("/login", func(c echo.Context) error {
		return c.String(http.StatusOK, "logged in")
	})

	s.Echo.POST("/logout", func(c echo.Context) error {
		return c.String(http.StatusOK, "logged out")
	})

	s.Echo.GET("/dashboard", func(c echo.Context) error {
		page := pages.Dashboard()
		return s.render(c, http.StatusOK, page)
	})

	s.Echo.GET("/complaints", func(c echo.Context) error {
		page := pages.Complaints()
		return s.render(c, http.StatusOK, page)
	})

	s.Echo.GET("/vendor-applications", func(c echo.Context) error {
		page := pages.Applications()
		return s.render(c, http.StatusOK, page)
	})

	s.Echo.GET("/verification-requests", func(c echo.Context) error {
		page := pages.IdentityVerification()
		return s.render(c, http.StatusOK, page)
	})

	s.Echo.GET("/account-management", func(c echo.Context) error {
		page := pages.AccountManagement()
		return s.render(c, http.StatusOK, page)
	})

	s.Echo.GET("/test", func(c echo.Context) error {
		page := pages.Experiment()
		return s.render(c, http.StatusOK, page)
	})
}
