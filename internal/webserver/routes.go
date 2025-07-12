package webserver

import (
	"bytes"
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	"github.com/project-agonyl/open-agonyl-servers/internal/webserver/mw"
	"github.com/rs/zerolog"
	"github.com/ziflex/lecho/v3"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Logger = lecho.From(s.logger.GetLoggerInstance().(zerolog.Logger))
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(mw.Logger(s.logger))
	csrfConfig := middleware.CSRFConfig{
		TokenLookup: "form:_csrf",
		ErrorHandler: func(err error, c echo.Context) error {
			return s.renderTemplate(c, "csrf-error", s.getBaseTemplateDataWithAuth(c), http.StatusForbidden)
		},
	}
	e.Use(middleware.CSRFWithConfig(csrfConfig))

	e.Renderer = NewTemplates()

	staticFS := GetStaticFS()
	assetsSubFS, err := fs.Sub(staticFS, "assets")
	if err != nil {
		s.logger.Error("Failed to create assets sub filesystem", shared.Field{Key: "error", Value: err})
	} else {
		e.StaticFS("/static", assetsSubFS)
	}

	e.Use(mw.AuthContextMiddleware(s.cfg, s.sessionStorage, s.logger))

	e.GET("/", s.handleIndex)

	e.GET("/about", s.handleAbout)

	e.GET("/register", s.handleRegisterPage)
	e.POST("/register", s.handleRegister)

	e.GET("/login", s.handleLoginPage)
	e.POST("/login", s.handleLogin)
	e.GET("/logout", s.handleLogout, mw.AuthGuardMiddleware(s.cfg, s.sessionStorage, s.logger))
	e.GET("/characters", s.handleCharacters, mw.AuthGuardMiddleware(s.cfg, s.sessionStorage, s.logger))

	return e
}

func (s *Server) getBaseTemplateData() map[string]interface{} {
	return map[string]interface{}{
		"ServerName": s.cfg.ServerName,
	}
}

func (s *Server) getBaseTemplateDataWithAuth(c echo.Context) map[string]interface{} {
	data := s.getBaseTemplateData()

	data["CSRFToken"] = c.Get(middleware.DefaultCSRFConfig.ContextKey)

	if mw.IsAuthenticated(c) {
		username, err := mw.GetUsername(c)
		if err == nil {
			data["IsAuthenticated"] = true
			data["Username"] = username
		} else {
			data["IsAuthenticated"] = false
		}
	} else {
		data["IsAuthenticated"] = false
	}

	return data
}

func (s *Server) renderTemplate(c echo.Context, templateName string, data interface{}, statusCode int) error {
	var buf bytes.Buffer
	if err := c.Echo().Renderer.Render(&buf, templateName, data, c); err != nil {
		return err
	}

	return c.HTML(statusCode, buf.String())
}
