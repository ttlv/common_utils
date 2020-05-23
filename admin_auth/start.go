package admin_auth

import (
	"fmt"
	"net/http"

	authboss "gopkg.in/authboss.v1"

	"github.com/jinzhu/configor"
	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/theplant/appkit/log"
)

var appConfig AppConfig
var DB *gorm.DB

type AppConfig struct {
	MountPath          string `required:"true"`
	CookieKey          string `required:"true"`
	CookieName         string `required:"true"`
	SessionKey         string `required:"true"`
	SuperAdminEmail    string `required:"true"`
	SuperAdminPassword string `required:"true"`
	LoginRedirectPath  string `required:"true"`
	LogoutRedirectPath string `required:"true"`
	ViewPath           string `required:"true"`
}
type ViewContext struct {
	Admin *admin.Admin
	Error string
}

// localesCollection
//    Possible Values: nil, []string, [][]string, func(interface{}, *qor.Context) [][]string, func(interface{}, *admin.Context) [][]string
//    Notes: If it is nil, will not show ViewableLocale and EditableLocale in User edit Page
func Start(l log.Logger, dbConn *gorm.DB, adm *admin.Admin, localesCollection interface{}) *authboss.Authboss {
	DB = dbConn
	DB.AutoMigrate(&AdminUser{})
	if err := configor.New(&configor.Config{ENVPrefix: "ADMIN_AUTH"}).Load(&appConfig); err != nil {
		l.Crit().Log("during", "configor.Load", "err", err, "msg", fmt.Sprintf("error loading config: %v", err))
	}

	user := AdminUser{}
	var adminUserCount int64
	err := DB.Model(&user).Count(&adminUserCount).Error
	if err != nil {
		panic(err)
	}
	DB.Where("email = ?", appConfig.SuperAdminEmail).Find(&user)
	if user.Role != "platform_admin" {
		user.Role = "platform_admin"
		user.Email = appConfig.SuperAdminEmail
		user.SetPassword(appConfig.SuperAdminPassword)
		user.Name = "平台管理员"
		DB.Save(&user)
	}

	setupAuth()
	adm.SetAuth(AdminAuth{})
	adm.RegisterViewPath("github.com/ttlv/common_utils/admin_auth/views")

	adm.GetRouter().Get("reset_password", func(c *admin.Context) {
		content := c.Render("reset_password", ViewContext{Admin: adm})
		c.Writer.Write([]byte(content))
	})
	adm.GetRouter().Post("reset_password", func(c *admin.Context) {
		u := c.CurrentUser.(*AdminUser)
		if !u.VaildPassword(c.Request.FormValue("old_password")) {
			content := c.Render("reset_password", ViewContext{
				Admin: adm,
				Error: "老密码不正确",
			})
			c.Writer.Write([]byte(content))
			return
		}
		if len(c.Request.FormValue("new_password")) < 8 {
			content := c.Render("reset_password", ViewContext{
				Admin: adm,
				Error: "密码不能低于8位",
			})
			c.Writer.Write([]byte(content))
			return
		}
		u.SetPassword(c.Request.FormValue("new_password"))
		c.DB.Save(&u)
		http.Redirect(c.Writer, c.Request, c.Admin.GetRouter().Prefix, 301)
	})
	return Auth
}
