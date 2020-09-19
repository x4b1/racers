package errors

import (
	"github.com/cockroachdb/errors"
)

func New(format string, args ...interface{}) error {
	return errors.Newf(format, args...)
}

func Wrap(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}
