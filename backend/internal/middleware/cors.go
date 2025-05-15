// internal/middleware/cors.go
package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// CORSConfig represents configuration for CORS middleware
type CORSConfig struct {
	AllowOrigins []string
	AllowMethods []string
}

// SetupCORS configures and returns CORS middleware for Echo
func SetupCORS(config CORSConfig) echo.MiddlewareFunc {
	// Set default values if not provided
	if len(config.AllowOrigins) == 0 {
		config.AllowOrigins = []string{"*"} // Allow all origins by default
	}
	
	if len(config.AllowMethods) == 0 {
		config.AllowMethods = []string{
			echo.GET, echo.HEAD, echo.PUT, echo.PATCH, 
			echo.POST, echo.DELETE, echo.OPTIONS,
		}
	}

	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     config.AllowOrigins,
		AllowMethods:     config.AllowMethods,
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		MaxAge:           86400, // Maximum value not ignored by any major browser
	})
}