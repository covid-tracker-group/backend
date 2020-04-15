package tokens

import (
	"time"

	"github.com/google/uuid"
	"simplon.biz/corona/pkg/config"
)

type TracingAuthenticationToken struct {
	BaseToken
}

func NewTracingAuthenticationToken() TracingAuthenticationToken {
	return TracingAuthenticationToken{
		BaseToken{
			Code:      uuid.New().String(),
			CreatedAt: time.Now(),
			ExpiresAt: time.Now().Add(config.ExpireDailyTracingTokensAfter),
		},
	}
}
