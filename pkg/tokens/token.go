package tokens

import (
	"time"
)

// Token is an token that can be used to submit or delete data
type Token interface {
	SetCode(string)
	GetCode() string
	GetCreatedAt() time.Time
}

// TokenCreator is a function that
type TokenCreator func(code string) Token
