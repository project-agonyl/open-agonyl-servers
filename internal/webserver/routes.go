package webserver

import (
	"bytes"
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
	"github.com/rs/zerolog"
	"github.com/ziflex/lecho/v3"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Logger = lecho.From(s.logger.GetLoggerInstance().(zerolog.Logger))
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CSRF())

	e.Renderer = NewTemplates()

	staticFS := GetStaticFS()
	assetsSubFS, err := fs.Sub(staticFS, "assets")
	if err != nil {
		s.logger.Error("Failed to create assets sub filesystem", shared.Field{Key: "error", Value: err})
	} else {
		e.StaticFS("/static", assetsSubFS)
	}

	e.GET("/", s.handleIndex)
	e.GET("/about", s.handleAbout)
	e.GET("/login", s.handleLoginPage)
	e.POST("/login", s.handleLogin)
	e.GET("/logout", s.handleLogout)

	e.GET("/characters", s.handleCharacters)

	return e
}

func (s *Server) getBaseTemplateData() map[string]interface{} {
	return map[string]interface{}{
		"ServerName": s.cfg.ServerName,
	}
}

func (s *Server) getBaseTemplateDataWithCSRF(c echo.Context) map[string]interface{} {
	data := s.getBaseTemplateData()
	data["CSRFToken"] = c.Get(middleware.DefaultCSRFConfig.ContextKey)
	return data
}

func (s *Server) renderTemplate(c echo.Context, templateName string, data interface{}, statusCode int) error {
	var buf bytes.Buffer
	if err := c.Echo().Renderer.Render(&buf, templateName, data, c); err != nil {
		return err
	}

	return c.HTML(statusCode, buf.String())
}
