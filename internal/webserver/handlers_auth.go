package webserver

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (s *Server) handleLoginPage(c echo.Context) error {
	data := s.getBaseTemplateDataWithCSRF(c)
	data["Title"] = "Login - " + s.cfg.ServerName
	return c.Render(http.StatusOK, "login", data)
}

func (s *Server) handleLogin(c echo.Context) error {
	username := strings.TrimSpace(c.FormValue("username"))
	password := strings.TrimSpace(c.FormValue("password"))

	if username == "" || password == "" {
		return s.renderTemplate(c, "login-error", map[string]interface{}{
			"Message": "Username and password are required.",
		}, http.StatusBadRequest)
	}

	if username == "player1" && password == "password123" {
		return s.renderTemplate(c, "login-success", s.getBaseTemplateData(), http.StatusOK)
	}

	return s.renderTemplate(c, "login-error", map[string]interface{}{
		"Message": "Invalid username or password. Please try again.",
	}, http.StatusUnauthorized)
}

func (s *Server) handleLogout(c echo.Context) error {
	return c.Redirect(http.StatusSeeOther, "/login")
}
