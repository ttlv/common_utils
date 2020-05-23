package admin_auth

import (
	"github.com/qor/admin"
	"github.com/qor/qor"
)

type AdminAuth struct {
}

func (AdminAuth) LoginURL(c *admin.Context) string {
	return "/auth/login"
}

func (AdminAuth) LogoutURL(c *admin.Context) string {
	return "/auth/logout"
}

func (AdminAuth) GetCurrentUser(c *admin.Context) qor.CurrentUser {
	userInter, err := Auth.CurrentUser(c.Writer, c.Request)
	if userInter != nil && err == nil {
		return userInter.(*AdminUser)
	}
	return nil
}
