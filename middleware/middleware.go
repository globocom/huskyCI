package middleware

import (
	"github.com/globocom/husky/config"
	"github.com/labstack/echo"
)

// RequestConfigMiddleware is a middleware to send request info to New Relic
func RequestConfigMiddleware(config *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("config", config)
			return next(c)
		}
	}
}
