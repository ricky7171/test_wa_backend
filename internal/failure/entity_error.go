package failure

import (
	"errors"
	"strings"
)

func ErrFieldRequired(nameFields ...string) error {
	var message string
	if len(nameFields) >= 1 {
		message = strings.Join(nameFields[:], " or ") + " cannot be empty "
	} else {
		message = "some field cannot be empty"
	}
	message = "301 - " + message
	return errors.New(message)
}

func ErrFieldLenConstraint(nameField string, minLength string, maxLength string) error {
	var message string
	message = nameField
	if minLength != "" {
		message += " should have min length " + minLength
		if maxLength != "" {
			message += "and "
		}
	}
	if maxLength != "" {
		message += "should have max length" + maxLength
	}

	message = "302 - " + message
	return errors.New(message)
}

func ErrFieldNumber(nameFields ...string) error {
	var message string
	if len(nameFields) == 0 {
		message = "some field should use number > -1"
	} else {
		message = strings.Join(nameFields, ",") + " should use number > -1"
	}
	message = "303 - " + message
	return errors.New(message)
}

func ErrPasswordNotMatch() error {
	return errors.New("305 - password doesn't match")
}
