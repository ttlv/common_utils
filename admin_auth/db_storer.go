package admin_auth

import "gopkg.in/authboss.v1"

type AuthStorer struct {
}

func (s AuthStorer) Create(key string, attr authboss.Attributes) error {
	var user AdminUser
	if err := attr.Bind(&user, true); err != nil {
		return err
	}

	if err := DB.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

func (s AuthStorer) Put(key string, attr authboss.Attributes) error {
	var user AdminUser
	if err := DB.Where("email = ?", key).First(&user).Error; err != nil {
		return authboss.ErrUserNotFound
	}

	if err := attr.Bind(&user, true); err != nil {
		return err
	}

	if err := DB.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

func (s AuthStorer) Get(key string) (result interface{}, err error) {
	var user AdminUser
	if err := DB.Where("email = ?", key).First(&user).Error; err != nil {
		return nil, authboss.ErrUserNotFound
	}
	return &user, nil
}

func (s AuthStorer) ConfirmUser(tok string) (result interface{}, err error) {
	var user AdminUser
	if err := DB.Where("confirm_token = ?", tok).First(&user).Error; err != nil {
		return nil, authboss.ErrUserNotFound
	}
	return &user, nil
}

func (s AuthStorer) RecoverUser(rec string) (result interface{}, err error) {
	var user AdminUser
	if err := DB.Where("recover_token = ?", rec).First(&user).Error; err != nil {
		return nil, authboss.ErrUserNotFound
	}
	return &user, nil
}
