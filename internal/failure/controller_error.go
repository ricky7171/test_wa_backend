package failure

import "errors"

func ErrLoginFailed() error {
	return errors.New("401 - phone or password is wrong")
}
