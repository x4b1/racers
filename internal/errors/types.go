package errors

import "encoding/json"

type internalError struct{ error }

func (e internalError) Unwrap() error { return e.error }

func NewInternalError(format string, args ...interface{}) error {
	return &internalError{New(format, args...)}
}

func WrapInternalError(err error) error {
	if err != nil {
		return &internalError{err}
	}

	return nil
}

func IsInternalError(err error) bool {
	var target *internalError
	return As(err, &target)
}

type notFoundError struct{ error }

func (e notFoundError) Unwrap() error { return e.error }

func WrapNotFoundError(err error) error {
	return &notFoundError{err}
}

func IsNotFoundError(err error) bool {
	var target *notFoundError
	return As(err, &target)
}

type wrongInputError struct{ error }

func (e wrongInputError) Unwrap() error { return e.error }

func WrapWrongInputError(err error) error {
	return &wrongInputError{err}
}

func IsWrongInputError(err error) bool {
	var target *wrongInputError
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
