package admin_auth

import (
	"net/http"
	"os"
	"regexp"

	"github.com/gorilla/csrf"
	_ "github.com/ttlv/common_utils/admin_auth/authboss/preauth"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/authboss.v1"
	_ "gopkg.in/authboss.v1/auth"
	_ "gopkg.in/authboss.v1/confirm"
	_ "gopkg.in/authboss.v1/recover"
	_ "gopkg.in/authboss.v1/register"
)

var Auth *authboss.Authboss

func setupAuth() {
	setupSessionStorer()
	setupCookieStorer()
	Auth = authboss.New()
	Auth.MountPath = appConfig.MountPath
	Auth.XSRFName = "gorilla.csrf.Token"
	Auth.XSRFMaker = func(_ http.ResponseWriter, r *http.Request) string {
		return csrf.Token(r)
	}
	Auth.CookieStoreMaker = NewCookieStorer
	Auth.SessionStoreMaker = NewSessionStorer
	Auth.BCryptCost = bcrypt.DefaultCost
	Auth.LogWriter = os.Stdout
	Auth.Storer = &AuthStorer{}
	Auth.ViewsPath = appConfig.ViewPath
	Auth.LayoutDataMaker = layoutData
	Auth.AuthLoginOKPath = appConfig.LoginRedirectPath
	Auth.AuthLogoutOKPath = appConfig.LogoutRedirectPath

	Auth.Policies = []authboss.Validator{
		authboss.Rules{
			FieldName:       "email",
			Required:        true,
			AllowWhitespace: false,
			MustMatch:       regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`),
			MatchError:      "Please input a valid email address",
		},
		authboss.Rules{
			FieldName:       "password",
			Required:        true,
			MinLength:       4,
			MaxLength:       8,
			AllowWhitespace: false,
		},
	}

	if err := Auth.Init("preauth"); err != nil {
		panic(err)
	}
}

func CurrentLocale(req *http.Request) string {
	locale := "en-US"
	if cookie, err := req.Cookie("locale"); err == nil {
		locale = cookie.Value
	}
	return locale
}

func layoutData(w http.ResponseWriter, r *http.Request) authboss.HTMLData {
	return authboss.HTMLData{}
}
