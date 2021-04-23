package failure

import (
	"errors"
)

func ErrRepoFailedQueryGet() error {
	return errors.New("001 - failed to query get data")
}

func ErrRepoFailedQueryInsert() error {
	return errors.New("002 - failed to query get insert")
}

func ErrRepoFailedQueryUpdate() error {
	return errors.New("003 - failed to query get update")
}

func ErrRepoFailedQueryDelete() error {
	return errors.New("004 - failed to query get delete")
}

func ErrRepoFailedConvert() error {
	return errors.New("005 - failed to convert given data")
}
