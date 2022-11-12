package google

import "github.com/pkg/errors"

var (
	CredentialIsNotExisted = errors.Errorf("Credential is not existed")
	TokenIsNotExisted      = errors.New("Token is not existed")
	TokenIsInvalid         = errors.New("Token is not valid")
)
