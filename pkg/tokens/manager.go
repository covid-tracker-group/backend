package tokens

import "time"

type TokenInfo struct {
	Token     string
	CreatedAt time.Time
}

// TokenManager is responsible for creating and validation authorisation tokens
type TokenManager interface {
	CreateToken() (string, error)
	VerifyToken(string) (bool, error)
	RetractToken(string) error
	Expire(onExpire func(string)) error
}
