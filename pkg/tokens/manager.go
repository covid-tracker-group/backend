package tokens

//go:generate mockgen -destination=../mocks/token_manager.go -package mocks . TokenManager

// TokenManager is responsible for creating and validation authorisation tokens
type TokenManager interface {
	StoreToken(Token) error
	VerifyToken(string) (bool, error)
	RetractToken(string) error
	Expire(onExpire func(string)) error
}
