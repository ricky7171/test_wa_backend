package failure

import "errors"

func ErrConvertObjectIdToHex() error {
	return errors.New("901 - failed to convert from object id to string")
}

func ErrHashPassword() error {
	return errors.New("902 - failed to hash password")
}
