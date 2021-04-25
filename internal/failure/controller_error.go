package failure

import "errors"

func ErrLoginFailed() error {
	return errors.New("401 - phone or password is wrong")
}

func ErrCannotReadJson() error {
	return errors.New("402 - cannot read JSON")
}
