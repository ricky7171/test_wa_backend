package failure

import "errors"

func ErrUserNotFound() error {
	return errors.New("201 - User not found")
}

func ErrDuplicatePhone() error {
	return errors.New("202 - this phone already exists")
}
