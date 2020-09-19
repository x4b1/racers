package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

func New(format string, args ...interface{}) error {
	return errors.New(fmt.Sprintf(format, args...))
}

func Wrap(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}
