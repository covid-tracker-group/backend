package tokens

// TokenManager is responsible for creating and validation authorisation tokens
type TokenManager interface {
	CreateToken() (string, error)
	VerifyToken(string) (bool, error)
	RetractToken(string) error
	Expire() error
}
