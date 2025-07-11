package mw

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/project-agonyl/open-agonyl-servers/internal/shared"
)

func Logger(logger shared.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			req := c.Request()
			res := c.Response()
			logger.Info("HTTP Request",
				shared.Field{Key: "method", Value: req.Method},
				shared.Field{Key: "uri", Value: req.RequestURI},
				shared.Field{Key: "status", Value: res.Status},
				shared.Field{Key: "latency", Value: time.Since(start).String()},
				shared.Field{Key: "remote_ip", Value: c.RealIP()},
				shared.Field{Key: "user_agent", Value: req.UserAgent()},
			)
			return err
		}
	}
}
