package webserver

import (
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/constants"
	"github.com/project-agonyl/open-agonyl-servers/internal/webserver/helpers"
	"github.com/project-agonyl/open-agonyl-servers/internal/webserver/mw"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) handleLoginPage(c echo.Context) error {
	data := s.getBaseTemplateDataWithAuth(c)
	data["Title"] = "Login - " + s.cfg.ServerName
	if mw.IsAuthenticated(c) {
		return c.Redirect(http.StatusSeeOther, "/")
	}

	return c.Render(http.StatusOK, "login", data)
}

func (s *Server) handleLogin(c echo.Context) error {
	username := strings.TrimSpace(c.FormValue("username"))
	password := strings.TrimSpace(c.FormValue("password"))
	account, err := s.db.GetAccountByUsername(username)
	if err != nil {
		return s.renderTemplate(c, "login-error", map[string]interface{}{
			"Message": "Invalid username or password. Please try again.",
		}, http.StatusUnauthorized)
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(password))
	if err != nil {
		return s.renderTemplate(c, "login-error", map[string]interface{}{
			"Message": "Invalid username or password. Please try again.",
		}, http.StatusUnauthorized)
	}

	if strings.EqualFold(account.Status, constants.AccountStatusBanned) {
		return s.renderTemplate(c, "login-error", map[string]interface{}{
			"Message": "Account is banned. Please contact support.",
		}, http.StatusUnauthorized)
	}

	if !strings.EqualFold(account.Status, constants.AccountStatusActive) {
		return s.renderTemplate(c, "login-error", map[string]interface{}{
			"Message": "Account is not active. Please activate your account.",
		}, http.StatusUnauthorized)
	}

	sessionID, err := s.sessionStorage.CreateSession(
		account.ID,
		c.Request().UserAgent(),
		c.RealIP(),
		time.Now().Add(time.Second*time.Duration(s.cfg.JwtExpiry)),
	)
	if err != nil {
		s.logger.Error("Failed to create session", shared.Field{Key: "error", Value: err})
		return s.renderTemplate(c, "login-error", map[string]interface{}{
			"Message": "Failed to login. Please try again later.",
		}, http.StatusInternalServerError)
	}

	claims := map[string]interface{}{
		"username":   account.Username,
		"session_id": sessionID,
		"sub":        account.AccountID,
		"aud":        s.cfg.ServerName,
	}
	accessToken, err := helpers.GenerateJwt(claims, s.cfg.JwtSecret, s.cfg.JwtExpiry)
	if err != nil {
		s.logger.Error("Failed to generate access token", shared.Field{Key: "error", Value: err})
		return s.renderTemplate(c, "login-error", map[string]interface{}{
			"Message": "Failed to login. Please try again later.",
		}, http.StatusInternalServerError)
	}

	helpers.WriteCookie(c, s.cfg.SessionCookieName, accessToken, s.cfg.JwtExpiry)
	return s.renderTemplate(c, "login-success", s.getBaseTemplateDataWithAuth(c), http.StatusOK)
}

func (s *Server) handleLogout(c echo.Context) error {
	if mw.IsAuthenticated(c) {
		sessionID, _ := mw.GetSessionID(c)
		if sessionID != "" {
			err := s.sessionStorage.RevokeSession(sessionID)
			if err != nil {
				s.logger.Error("Failed to revoke session", shared.Field{Key: "error", Value: err})
			}
		}
	}

	helpers.DeleteCookie(c, s.cfg.SessionCookieName)
	return c.Redirect(http.StatusSeeOther, "/login")
}
