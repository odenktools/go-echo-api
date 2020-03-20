package oauth

import (
	"github.com/labstack/echo"
)

const (
	tokensResource     = "token"
	tokensPath         = "/" + tokensResource
	introspectResource = "introspect"
	introspectPath     = "/" + introspectResource
)

// RegisterRoutes registers route handlers for the oauth service
func (s *Service) RegisterRoutes(e *echo.Echo, path string) {
	auth := e.Group(path)
	auth.POST(tokensPath, s.oauthToken)
}