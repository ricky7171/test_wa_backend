package failure

import (
	"errors"
)

func ErrGenerateToken() error {
	return errors.New("101 - failed to generate token")
}

func ErrFailedParseClaim() error {
	return errors.New("102 - failed to parse claims token")
}

func ErrTokenInvalid() error {
	return errors.New("103 - token is invalid")
}

func ErrTokenExpired() error {
	return errors.New("104 - token is expired")
}
