package services

import (
	"reflect"

	"go-echo-api/config"
	"go-echo-api/oauth"
	"go-echo-api/session"

	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
)

func init() {

}

var (

	// OauthService ...
	OauthService oauth.ServiceInterface

	// SessionService ...
	SessionService session.ServiceInterface
)

// UseOauthService sets the oAuth service
func UseOauthService(o oauth.ServiceInterface) {
	OauthService = o
}

// UseSessionService sets the session service
func UseSessionService(s session.ServiceInterface) {
	SessionService = s
}

// Init starts up all services
func Init(cnf *config.Config, db *gorm.DB) error {

	if nil == reflect.TypeOf(OauthService) {
		OauthService = oauth.NewService(cnf, db)
	}

	if nil == reflect.TypeOf(SessionService) {
		// note: default session store is CookieStore
		SessionService = session.NewService(cnf, sessions.NewCookieStore([]byte(cnf.Session.Secret)))
	}

	return nil
}

// Close closes any open services
func Close() {
	OauthService.Close()
	SessionService.Close()
}
