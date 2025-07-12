package webserver

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) handleCharacters(c echo.Context) error {
	data := s.getBaseTemplateDataWithAuth(c)
	data["Title"] = "Characters - " + s.cfg.ServerName
	return c.Render(http.StatusOK, "characters", data)
}
