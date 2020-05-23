package preauth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ttlv/common_utils/admin_auth/authboss/response"
	"github.com/ttlv/common_utils/services/sms"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/authboss.v1"
)

var SmsService sms.SmsServerInterface
var Rander RandMaker

const (
	methodGET       = "GET"
	methodPOST      = "POST"
	tplLogin        = "login.html.tpl"
	tplTwoSteplogin = "twostep_login.html.tpl"
)

func init() {
	authboss.RegisterModule("preauth", &PreAuth{})
	SmsService, _ = sms.NewSmsService()
	Rander = DefaultRandMaker{}
}

// Auth module
type PreAuth struct {
	*authboss.Authboss
	templates response.Templates
}

// Initialize module
func (a *PreAuth) Initialize(ab *authboss.Authboss) (err error) {
	a.Authboss = ab

	if a.Storer == nil && a.StoreMaker == nil {
		return errors.New("auth: Need a Storer")
	}

	if len(a.XSRFName) == 0 {
		return errors.New("auth: XSRFName must be set")
	}

	if a.XSRFMaker == nil {
		return errors.New("auth: XSRFMaker must be defined")
	}

	a.templates, err = response.LoadTemplates(a.Authboss, a.Layout, a.ViewsPath, tplLogin, tplTwoSteplogin)
	if err != nil {
		return err
	}
	return nil
}

// Routes for the module
func (a *PreAuth) Routes() authboss.RouteTable {
	return authboss.RouteTable{
		"/preauth": a.PreAuthHandlerFunc,
		"/login":   a.loginHandlerFunc,
		"/logout":  a.logoutHandlerFunc,
	}
}

// Storage requirements
func (a *PreAuth) Storage() authboss.StorageOptions {
	return authboss.StorageOptions{
		a.PrimaryID:            authboss.String,
		authboss.StorePassword: authboss.String,
	}
}

func (a *PreAuth) PreAuthHandlerFunc(ctx *authboss.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case methodPOST:
		key := r.FormValue(a.PrimaryID)
		password := r.FormValue("password")
		data := authboss.NewHTMLData(
			"primaryID", a.PrimaryID,
			"primaryIDValue", key,
		)

		valid, err := validateCredentials(ctx, key, password)
		if err != nil || !valid {
			data["error"] = "邮箱或密码不正确."
			return a.templates.Render(ctx, w, r, tplLogin, data)
		}

		user, err := a.Authboss.Storer.Get(key)
		if err != nil {
			data["error"] = "找不到用户."
			return a.templates.Render(ctx, w, r, tplLogin, data)
		}
		phone := strings.TrimSpace(user.(UserInterface).GetPhoneNumber())
		if phone == "" {
			data["error"] = "用户的手机号码不能为空."
			return a.templates.Render(ctx, w, r, tplLogin, data)
		}

		attrs := make(authboss.Attributes)
		attrs["opt"] = fmt.Sprintf("%v", Rander.Gen(6))
		attrs["expire_at"] = time.Now().Add(300 * time.Second)
		attrs["token"], _ = generateRandomString(32)
		if err := a.Authboss.Storer.Put(key, attrs); err != nil {
			w.Write([]byte(`{ "error": "不能获取到Opt." }`))
			return nil
		}

		data["token"] = attrs["token"]
		data["phone"] = hiddenPhone(phone)
		SmsService.Send("AdminAuth", "", phone, fmt.Sprintf("【Admin】登陆验证码: %v", attrs["opt"]))
		return a.templates.Render(ctx, w, r, tplTwoSteplogin, data)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	return nil
}

func (a *PreAuth) loginHandlerFunc(ctx *authboss.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case methodGET:
		data := authboss.NewHTMLData(
			"showRemember", a.IsLoaded("remember"),
			"showRecover", a.IsLoaded("recover"),
			"showRegister", a.IsLoaded("register"),
			"primaryID", a.PrimaryID,
			"primaryIDValue", "",
		)
		return a.templates.Render(ctx, w, r, tplLogin, data)
	case methodPOST:
		key := r.FormValue(a.PrimaryID)

		errData := authboss.NewHTMLData(
			"primaryID", a.PrimaryID,
			"primaryIDValue", key,
		)

		if r.FormValue("token") == "" || r.FormValue("code") == "" || key == "" {
			errData["error"] = "邮箱, Token和验证码不能为空."
			return a.templates.Render(ctx, w, r, tplLogin, errData)
		}

		user, err := a.Authboss.Storer.Get(key)
		if user.(UserInterface).GetToken() != r.FormValue("token") {
			errData["error"] = "Token不一致"
			return a.templates.Render(ctx, w, r, tplLogin, errData)
		}
		if user.(UserInterface).GetOtp() != r.FormValue("code") {
			errData["error"] = "验证码不正确"
			return a.templates.Render(ctx, w, r, tplLogin, errData)
		}
		if time.Now().After(user.(UserInterface).GetExpireAt()) {
			errData["error"] = "验证码已经过期, 请重新登陆"
			return a.templates.Render(ctx, w, r, tplLogin, errData)
		}

		interrupted, err := a.Callbacks.FireBefore(authboss.EventAuth, ctx)
		if err != nil {
			return err
		} else if interrupted != authboss.InterruptNone {
			var reason string
			switch interrupted {
			case authboss.InterruptAccountLocked:
				reason = "Your account has been locked."
			case authboss.InterruptAccountNotConfirmed:
				reason = "Your account has not been confirmed."
			}
			response.Redirect(ctx, w, r, a.AuthLoginFailPath, "", reason, false)
			return nil
		}

		ctx.SessionStorer.Put(authboss.SessionKey, key)
		ctx.SessionStorer.Del(authboss.SessionHalfAuthKey)
		ctx.Values = map[string]string{authboss.CookieRemember: r.FormValue(authboss.CookieRemember)}

		if err := a.Callbacks.FireAfter(authboss.EventAuth, ctx); err != nil {
			return err
		}
		response.Redirect(ctx, w, r, a.AuthLoginOKPath, "", "", true)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	return nil
}

func (a *PreAuth) logoutHandlerFunc(ctx *authboss.Context, w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case methodGET:
		ctx.SessionStorer.Del(authboss.SessionKey)
		ctx.CookieStorer.Del(authboss.CookieRemember)
		ctx.SessionStorer.Del(authboss.SessionLastAction)

		response.Redirect(ctx, w, r, a.AuthLogoutOKPath, "You have logged out", "", true)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	return nil
}

func validateCredentials(ctx *authboss.Context, key, password string) (bool, error) {
	if err := ctx.LoadUser(key); err == authboss.ErrUserNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}

	actualPassword, err := ctx.User.StringErr(authboss.StorePassword)
	if err != nil {
		return false, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(actualPassword), []byte(password)); err != nil {
		return false, nil
	}

	return true, nil
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}
	return b, nil
}

func generateRandomString(s int) (string, error) {
	b, err := generateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func hiddenPhone(phone string) string {
	slice := strings.Split(phone, "")
	str := strings.Join(slice[0:5], "") + "****" + strings.Join(slice[len(phone)-3:len(phone)], "")
	return str
}
