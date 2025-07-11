package webserver

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) handleIndex(c echo.Context) error {
	data := s.getBaseTemplateData()
	data["Title"] = "Home - " + s.cfg.ServerName
	return c.Render(http.StatusOK, "index", data)
}

func (s *Server) handleAbout(c echo.Context) error {
	data := s.getBaseTemplateData()
	data["Title"] = "About - " + s.cfg.ServerName
	return c.Render(http.StatusOK, "about", data)
}
