package authz

type MedicalAuthCodeError struct {
	msg string
}

type MedialAuthInfo struct{}

func (mae MedicalAuthCodeError) Error() string {
	return mae.msg
}

type AuthorisationManager struct{}

func NewAuthorisationManager() *AuthorisationManager {
	return &AuthorisationManager{}
}

func (mgr *AuthorisationManager) ValidateMedicalAuthCode(code string) (*MedialAuthInfo, error) {
	return &MedialAuthInfo{}, nil
}
