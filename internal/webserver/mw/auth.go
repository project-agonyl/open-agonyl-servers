package mw

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	"github.com/project-agonyl/open-agonyl-servers/internal/webserver/config"
	"github.com/project-agonyl/open-agonyl-servers/internal/webserver/helpers"
	"github.com/project-agonyl/open-agonyl-servers/internal/webserver/session"
)

const (
	ContextKeySessionID   = "session_id"
	ContextKeyAccountID   = "account_id"
	ContextKeyUsername    = "username"
	ContextKeySession     = "session"
	ContextKeyAccessToken = "access_token"
)

func AuthGuardMiddleware(cfg *config.EnvVars, sessionStorage session.SessionStorage, logger shared.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if IsAuthenticated(c) {
				return next(c)
			}

			accessToken, err := helpers.ReadCookie(c, cfg.SessionCookieName)
			if err != nil {
				logger.Debug("No session cookie found", shared.Field{Key: "error", Value: err})
				return c.Redirect(http.StatusSeeOther, "/login")
			}

			claims, err := helpers.VerifyAndParseJwt(accessToken, cfg.JwtSecret)
			if err != nil {
				logger.Debug("Invalid JWT token", shared.Field{Key: "error", Value: err})
				return c.Redirect(http.StatusSeeOther, "/login")
			}

			session, err := sessionStorage.GetSessionBySessionID(claims.SessionID, true)
			if err != nil || session == nil {
				logger.Debug("Session not found or inactive", shared.Field{Key: "session_id", Value: claims.SessionID})
				return c.Redirect(http.StatusSeeOther, "/login")
			}

			c.Set(ContextKeySession, session)
			c.Set(ContextKeyAccessToken, accessToken)
			c.Set(ContextKeyAccountID, claims.Subject)
			c.Set(ContextKeyUsername, claims.Username)
			c.Set(ContextKeySessionID, claims.SessionID)
			return next(c)
		}
	}
}

func AuthContextMiddleware(cfg *config.EnvVars, sessionStorage session.SessionStorage, logger shared.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			accessToken, err := helpers.ReadCookie(c, cfg.SessionCookieName)
			if err != nil {
				return next(c)
			}

			claims, err := helpers.VerifyAndParseJwt(accessToken, cfg.JwtSecret)
			if err != nil {
				return next(c)
			}

			session, err := sessionStorage.GetSessionBySessionID(claims.SessionID, true)
			if err != nil || session == nil {
				return next(c)
			}

			c.Set(ContextKeySession, session)
			c.Set(ContextKeyAccessToken, accessToken)
			c.Set(ContextKeyAccountID, claims.Subject)
			c.Set(ContextKeyUsername, claims.Username)
			c.Set(ContextKeySessionID, claims.SessionID)
			return next(c)
		}
	}
}

func IsAuthenticated(c echo.Context) bool {
	sessionID, ok := c.Get(ContextKeySessionID).(string)
	return ok && sessionID != ""
}

func GetSessionID(c echo.Context) (string, error) {
	sessionID, ok := c.Get(ContextKeySessionID).(string)
	if !ok {
		return "", errors.New("session_id not found in context")
	}

	return sessionID, nil
}

func GetAccountID(c echo.Context) (uint32, error) {
	accountID, ok := c.Get(ContextKeyAccountID).(uint32)
	if !ok {
		return 0, errors.New("account_id not found in context")
	}

	return accountID, nil
}

func GetUsername(c echo.Context) (string, error) {
	username, ok := c.Get(ContextKeyUsername).(string)
	if !ok {
		return "", errors.New("username not found in context")
	}

	return username, nil
}

func GetSession(c echo.Context) (*session.Session, error) {
	session, ok := c.Get(ContextKeySession).(*session.Session)
	if !ok {
		return nil, errors.New("session not found in context")
	}

	return session, nil
}

func GetAccessToken(c echo.Context) (string, error) {
	accessToken, ok := c.Get(ContextKeyAccessToken).(string)
	if !ok {
		return "", errors.New("access_token not found in context")
	}

	return accessToken, nil
}
