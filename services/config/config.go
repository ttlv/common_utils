package config

import (
	"github.com/jinzhu/configor"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type CoinxServiceConfig struct {
	ServerAddr string
}

var _CoinxServiceConfig *CoinxServiceConfig

type QuotableServiceConfig struct {
	ServerAddr string
}

var _QuotableServiceConfig *QuotableServiceConfig

func MustGetCoinxServiceConfig() CoinxServiceConfig {
	if _CoinxServiceConfig != nil {
		return *_CoinxServiceConfig
	}

	coinxServiceConfig := &CoinxServiceConfig{}
	err := configor.New(&configor.Config{ENVPrefix: "COINX_SERVICE"}).Load(coinxServiceConfig)
	if err != nil {
		panic(err)
	}
	_CoinxServiceConfig = coinxServiceConfig

	return *_CoinxServiceConfig
}

type NotifierServiceConfig struct {
	ServerAddr string
}

var _NotifierServiceConfig *NotifierServiceConfig

func MustGetNotifierServiceConfig() NotifierServiceConfig {
	if _NotifierServiceConfig != nil {
		return *_NotifierServiceConfig
	}

	notifierServiceConfig := &NotifierServiceConfig{}
	err := configor.New(&configor.Config{ENVPrefix: "NOTIFIER_SERVICE"}).Load(notifierServiceConfig)
	if err != nil {
		panic(err)
	}
	_NotifierServiceConfig = notifierServiceConfig

	return *_NotifierServiceConfig
}

func MustGetQuotableServiceConfig() QuotableServiceConfig {
	if _QuotableServiceConfig != nil {
		return *_QuotableServiceConfig
	}

	coinxServiceConfig := &QuotableServiceConfig{}
	err := configor.New(&configor.Config{ENVPrefix: "QUOTABLE_SERVICE"}).Load(coinxServiceConfig)
	if err != nil {
		panic(err)
	}
	_QuotableServiceConfig = coinxServiceConfig

	return *_QuotableServiceConfig
}

type SmsServiceConfig struct {
	ServerAddr string
}

var _SmsServiceConfig *SmsServiceConfig

func MustGetSmsServiceConfig() SmsServiceConfig {
	if _SmsServiceConfig != nil {
		return *_SmsServiceConfig
	}

	smsServiceConfig := &SmsServiceConfig{}
	err := configor.New(&configor.Config{ENVPrefix: "SMS_SERVICE"}).Load(smsServiceConfig)
	if err != nil {
		panic(err)
	}
	_SmsServiceConfig = smsServiceConfig

	return *_SmsServiceConfig
}

type EmailServiceConfig struct {
	ServerAddr string
}

var _EmailServiceConfig *EmailServiceConfig

func MustGetEmailServiceConfig() EmailServiceConfig {
	if _EmailServiceConfig != nil {
		return *_EmailServiceConfig
	}

	emailServiceConfig := &EmailServiceConfig{}
	err := configor.New(&configor.Config{ENVPrefix: "EMAIL_SERVICE"}).Load(emailServiceConfig)
	if err != nil {
		panic(err)
	}
	_EmailServiceConfig = emailServiceConfig

	return *_EmailServiceConfig
}
