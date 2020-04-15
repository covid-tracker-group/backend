package tokens

import (
	"time"

	"simplon.biz/corona/pkg/config"
	"simplon.biz/corona/pkg/tools"
)

type HealthTestAuthenticationToken struct {
	BaseToken
	// Uid is the uid as returned in the SAML assertion from SIAM
	Uid string `json:"uid"`
}

func NewHealthTestAuthenticationToken(uid string) HealthTestAuthenticationToken {
	return HealthTestAuthenticationToken{
		BaseToken: BaseToken{
			Code:      tools.GenerateCode(),
			CreatedAt: time.Now(),
			ExpiresAt: time.Now().Add(config.ExpireDailyTracingTokensAfter).Round(time.Hour),
		},
		Uid: uid,
	}
}
