package webserver

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared/constants"
	"github.com/project-agonyl/open-agonyl-servers/internal/webserver/helpers"
	"github.com/project-agonyl/open-agonyl-servers/internal/webserver/mw"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Username       string `form:"username" validate:"required,min=6,max=20,username"`
	Email          string `form:"email" validate:"required,email"`
	Password       string `form:"password" validate:"required,min=6,max=20"`
	RepeatPassword string `form:"repeat_password" validate:"required,eqfield=Password"`
}

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

func (s *Server) handleRegisterPage(c echo.Context) error {
	data := s.getBaseTemplateDataWithAuth(c)
	data["Title"] = "Register - " + s.cfg.ServerName
	if mw.IsAuthenticated(c) {
		return c.Redirect(http.StatusSeeOther, "/")
	}

	return c.Render(http.StatusOK, "register", data)
}

func (s *Server) handleRegister(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return s.renderTemplate(c, "register-error", map[string]interface{}{
			"Message": "Invalid form data.",
		}, http.StatusBadRequest)
	}

	validate := validator.New()
	_ = validate.RegisterValidation("username", helpers.ValidateUsername)
	err := validate.Struct(req)
	if err != nil {
		fieldErrors := make(map[string]string)
		for _, ferr := range err.(validator.ValidationErrors) {
			field := strings.ToLower(ferr.Field())
			switch field {
			case "username":
				fieldErrors[field] = "Username must be between 6 and 20 characters, lowercase letters and numbers only."
			case "email":
				fieldErrors[field] = "Please enter a valid email address."
			case "password":
				fieldErrors[field] = "Password must be between 6 and 20 characters."
			case "repeatpassword":
				fieldErrors[field] = "Passwords do not match."
			default:
				fieldErrors[field] = "Invalid value."
			}
		}
		return c.Render(http.StatusBadRequest, "register-error", map[string]interface{}{
			"FieldErrors": fieldErrors,
		})
	}

	if acc, _ := s.db.GetAccountByUsername(req.Username); acc != nil {
		return c.Render(http.StatusBadRequest, "register-error", map[string]interface{}{
			"FieldErrors": map[string]string{"username": "Username is already taken."},
		})
	}

	if acc, _ := s.db.GetAccountByEmail(req.Email); acc != nil {
		return c.Render(http.StatusBadRequest, "register-error", map[string]interface{}{
			"FieldErrors": map[string]string{"email": "Email is already registered."},
		})
	}

	isVerification := s.cfg.IsAccountVerificationRequired
	_, err = s.db.CreateAccount(req.Username, req.Password, req.Email, isVerification)
	if err != nil {
		s.logger.Error("Failed to create account", shared.Field{Key: "error", Value: err})
		return s.renderTemplate(c, "register-error", map[string]interface{}{
			"Message": "Failed to create account. Please try again later.",
		}, http.StatusInternalServerError)
	}

	msg := "Account created successfully! Have fun in the game!"
	if isVerification {
		msg = "Account created! Please check your email to verify your account."
	}

	return s.renderTemplate(c, "register-success", map[string]interface{}{"Message": msg}, http.StatusOK)
}
