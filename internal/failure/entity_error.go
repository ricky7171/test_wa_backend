package failure

import (
	"errors"
	"strings"
)

func ErrFieldRequired(nameFields ...string) error {
	if len(nameFields) >= 1 {
		return errors.New(strings.Join(nameFields[:], " or ") + " cannot be empty ")
	}
	return errors.New("some field cannot be empty")
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

	message += "301 - "
	return errors.New(message)
}

func ErrFieldNumber(nameFields ...string) error {
	var message string
	if len(nameFields) == 0 {
		message = "some field should use number > -1"
	} else {
		message = strings.Join(nameFields, ",") + " should use number > -1"
	}
	message += "301 - "
	return errors.New(message)
}
