package admin_auth

import (
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"fmt"

	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
)

type AdminUser struct {
	gorm.Model
	Email         string
	Password      string
	Name          string
	RoleID        uint
	Role          string
	ViewDashboard bool
	ViewQuotable  bool
	// Confirm
	ConfirmToken string
	Confirmed    bool

	// Recover
	RecoverToken       string
	RecoverTokenExpiry *time.Time
	ViewableLocale     string
	EditableLocale     string
	ExternalID         uint

	PhoneNumber string
	Opt         string
	ExpireAt    time.Time
	Token       string
}

func (u AdminUser) ViewableLocales() []string {
	return strings.Split(u.ViewableLocale, ",")
}

func (u AdminUser) EditableLocales() []string {
	return strings.Split(u.EditableLocale, ",")
}

func (user AdminUser) DisplayName() string {
	return user.Name
}

func (user AdminUser) GetExternalID() uint {
	return user.ExternalID
}

func (user AdminUser) GetRoleID() uint {
	return user.RoleID
}

func (user AdminUser) GetRole() string {
	return user.Role
}

func (user AdminUser) GetId() string {
	return fmt.Sprintf("%v", user.ID)
}

func (user AdminUser) GetPhoneNumber() string {
	return user.PhoneNumber
}

func (user AdminUser) GetExpireAt() time.Time {
	return user.ExpireAt
}

func (user AdminUser) GetToken() string {
	return user.Token
}

func (user AdminUser) GetOtp() string {
	return user.Opt
}

func (user *AdminUser) SetPassword(password string) {
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		color.YellowString("Can't set password, get (%v)", err)
	}
	user.Password = string(bcryptPassword)
}

func (user *AdminUser) VaildPassword(password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err == nil {
		return true
	}
	return false
}

func (user *AdminUser) BeforeCreate() (err error) {
	user.Confirmed = true
	return
}
