package admin_auth

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"gopkg.in/authboss.v1"
)

var sessionCookieName string
var sessionStore *sessions.CookieStore

type SessionStorer struct {
	w http.ResponseWriter
	r *http.Request
}

func setupSessionStorer() {
	sessionCookieName = appConfig.CookieName
	sessionStoreKey, err := base64.StdEncoding.DecodeString(appConfig.SessionKey)
	if err != nil {
		panic(fmt.Sprintf("AdminAuth: init cookie store, got (%v)", err))
	}
	sessionStore = sessions.NewCookieStore(sessionStoreKey)
}

func NewSessionStorer(w http.ResponseWriter, r *http.Request) authboss.ClientStorer {
	return &SessionStorer{w, r}
}

func (s SessionStorer) Get(key string) (string, bool) {
	sessionStore.Options.HttpOnly = true
	session, err := sessionStore.Get(s.r, sessionCookieName)
	if err != nil {
		fmt.Println(err)
		return "", false
	}

	strInf, ok := session.Values[key]
	if !ok {
		return "", false
	}

	str, ok := strInf.(string)
	if !ok {
		return "", false
	}

	return str, true
}

func (s SessionStorer) Put(key, value string) {
	session, err := sessionStore.Get(s.r, sessionCookieName)
	if err != nil {
		fmt.Println(err)
		return
	}

	session.Values[key] = value
	session.Save(s.r, s.w)
}

func (s SessionStorer) Del(key string) {
	session, err := sessionStore.Get(s.r, sessionCookieName)
	if err != nil {
		fmt.Println(err)
		return
	}

	delete(session.Values, key)
	session.Save(s.r, s.w)
}
