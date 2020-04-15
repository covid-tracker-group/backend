package tokens

// TokenManager is responsible for creating and validation authorisation tokens
type TokenManager interface {
	StoreToken(Token) error
	VerifyToken(string) (bool, error)
	RetractToken(string) error
	Expire(onExpire func(string)) error
}
