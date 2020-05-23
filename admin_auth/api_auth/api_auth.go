package api_auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	authboss "gopkg.in/authboss.v1"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/tidwall/gjson"
)

type Claims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

type APIAuth struct {
	getTokenPath string
	secretKey    string
	ab           *authboss.Authboss
}

type ErrorResp struct {
	Error ErrorData `json:"error"`
}

type SuccessResp struct {
	Data SuccessData `json:"data"`
}

type SuccessData struct {
	Token string `json:"token"`
}

type ErrorData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type MatcherFunc func(url string, claims Claims, r *http.Request) bool

func New(ab *authboss.Authboss, secretKey string) APIAuth {
	return APIAuth{secretKey: secretKey, ab: ab}
}

func (auth APIAuth) HandleFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		(w).Header().Set("Access-Control-Allow-Origin", "*")
		(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		mySigningKey := []byte(auth.secretKey)

		body, _ := ioutil.ReadAll(r.Body)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		params := gjson.Parse(string(body))
		email := params.Get("email").Str
		password := params.Get("password").Str

		if email == "" && password == "" {
			email = r.FormValue("email")
			password = r.FormValue("password")
		}

		ctx := auth.ab.NewContext()
		if _, err, errResp := validateCredentials(ctx, email, password); err != nil {
			errResp, _ := json.Marshal(errResp)
			w.Write(errResp)
			return
		}

		claims := Claims{
			Email: ctx.User["email"].(string),
			Role:  ctx.User["role"].(string),
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: 24 * 60 * 60 * int64(time.Second),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		ss, _ := token.SignedString(mySigningKey)
		successJSON, _ := json.Marshal(SuccessResp{Data: SuccessData{Token: ss}})
		w.Write(successJSON)
	}
}

func (auth APIAuth) Middleware(matcherFunc MatcherFunc) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var (
				authorization = r.Header.Get("Authorization")
				claim         = Claims{}
			)
			if strings.HasPrefix(authorization, "Bearer ") {
				tokenStr := strings.Replace(authorization, "Bearer ", "", -1)
				jwt.ParseWithClaims(tokenStr, &claim, func(token *jwt.Token) (interface{}, error) {
					return []byte(auth.secretKey), nil
				})
			}
			if matcherFunc(r.URL.String(), claim, r) {
				handler.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				errorJSON, _ := json.Marshal(ErrorResp{Error: ErrorData{Code: http.StatusUnauthorized, Message: "User Unauthorized"}})
				w.Write(errorJSON)
			}
			return
		})
	}
}

func validateCredentials(ctx *authboss.Context, key, password string) (bool, error, ErrorResp) {
	if err := ctx.LoadUser(key); err != nil {
		return false, fmt.Errorf("User not found."), ErrorResp{Error: ErrorData{Code: 404, Message: "User not found."}}
	}

	actualPassword, err := ctx.User.StringErr(authboss.StorePassword)
	if err != nil {
		return false, err, ErrorResp{Error: ErrorData{Code: 400, Message: err.Error()}}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(actualPassword), []byte(password)); err != nil {
		return false, fmt.Errorf("Password is incorrect"), ErrorResp{Error: ErrorData{Code: 401, Message: "Password is incorrect"}}
	}

	return true, nil, ErrorResp{}
}
