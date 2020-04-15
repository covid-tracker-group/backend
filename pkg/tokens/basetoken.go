package tokens

import (
	"time"
)

type BaseToken struct {
	Code      string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (t BaseToken) GetCode() string {
	return t.Code
}

func (t BaseToken) SetCode(code string) {
	t.Code = code
}

func (t BaseToken) GetCreatedAt() time.Time {
	return t.CreatedAt
}

func (t BaseToken) GetExpiresAt() time.Time {
	return t.ExpiresAt
}
