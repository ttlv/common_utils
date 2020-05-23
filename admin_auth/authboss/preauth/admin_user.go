package preauth

import "time"

type UserInterface interface {
	GetPhoneNumber() string
	GetExpireAt() time.Time
	GetToken() string
	GetOtp() string
}
