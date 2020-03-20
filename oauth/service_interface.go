package oauth

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"

	"go-echo-api/config"
	"go-echo-api/entity"
	"go-echo-api/session"
	"go-echo-api/utils/routes"
)

// ServiceInterface defines exported methods
type ServiceInterface interface {
	// Exported methods
	GetConfig() *config.Config
	RestrictToRoles(allowedRoles ...string)
	IsRoleAllowed(role string) bool
	FindRoleByID(id string) (*entity.OauthRole, error)
	GetRoutes() []routes.Route
	RegisterRoutes(router *echo.Router, prefix string)
	ClientExists(clientID string) bool
	FindClientByClientID(clientID string) (*entity.OauthClient, error)
	CreateClient(clientID, secret, redirectURI string) (*entity.OauthClient, error)
	CreateClientTx(tx *gorm.DB, clientID, secret, redirectURI string) (*entity.OauthClient, error)
	AuthClient(clientID, secret string) (*entity.OauthClient, error)
	UserExists(username string) bool
	FindUserByUsername(username string) (*entity.OauthUser, error)
	CreateUser(roleID, username, password string) (*entity.OauthUser, error)
	CreateUserTx(tx *gorm.DB, roleID, username, password string) (*entity.OauthUser, error)
	SetPassword(user *entity.OauthUser, password string) error
	SetPasswordTx(tx *gorm.DB, user *entity.OauthUser, password string) error
	UpdateUsername(user *entity.OauthUser, username string) error
	UpdateUsernameTx(db *gorm.DB, user *entity.OauthUser, username string) error
	AuthUser(username, thePassword string) (*entity.OauthUser, error)
	GetScope(requestedScope string) (string, error)
	GetDefaultScope() string
	ScopeExists(requestedScope string) bool
	Login(client *entity.OauthClient, user *entity.OauthUser, scope string) (*entity.OauthAccessToken, *entity.OauthRefreshToken, error)
	GrantAuthorizationCode(client *entity.OauthClient, user *entity.OauthUser, expiresIn int, redirectURI, scope string) (*entity.OauthAuthorizationCode, error)
	GrantAccessToken(client *entity.OauthClient, user *entity.OauthUser, expiresIn int, scope string) (*entity.OauthAccessToken, error)
	GetOrCreateRefreshToken(client *entity.OauthClient, user *entity.OauthUser, expiresIn int, scope string) (*entity.OauthRefreshToken, error)
	GetValidRefreshToken(token string, client *entity.OauthClient) (*entity.OauthRefreshToken, error)
	Authenticate(token string) (*entity.OauthAccessToken, error)
	NewIntrospectResponseFromAccessToken(accessToken *entity.OauthAccessToken) (*IntrospectResponse, error)
	NewIntrospectResponseFromRefreshToken(refreshToken *entity.OauthRefreshToken) (*IntrospectResponse, error)
	ClearUserTokens(userSession *session.UserSession)
	Close()
}
