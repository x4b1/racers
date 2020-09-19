package errors

import "encoding/json"

type internalError struct{ error }

func (e internalError) Unwrap() error { return e.error }

func WrapInternalError(err error, format string, args ...interface{}) error {
	return &internalError{Wrap(err, format, args...)}
}

func IsInternalError(err error) bool {
	var target *internalError
	return As(err, &target)
}

type notFoundError struct{ error }

func (e notFoundError) Unwrap() error { return e.error }

func WrapNotFoundError(err error, format string, args ...interface{}) error {
	return &notFoundError{Wrap(err, format, args...)}
}

func IsNotFoundError(err error) bool {
	var target *notFoundError
	return As(err, &target)
}

type ValidationError struct {
	Errors []error
}

func (ve *ValidationError) Add(err error) {
	ve.Errors = append(ve.Errors, err)
}

func (ve ValidationError) Error() string {
	b, _ := json.Marshal(ve.Errors)

	return string(b)
}

func (ve ValidationError) Valid() error {
	if len(ve.Errors) == 0 {
		return nil
	}

	return ve
}
