package errs

import (
	"net/http"

	"github.com/pkg/errors"
)

type ErrorType uint

const (
	Unknown ErrorType = iota
	Invalidated
	Unauthorized
	Forbidden
	NotFound
	Conflict
	Failed
)

type typeGetter interface {
	Type() ErrorType
}

type customError struct {
	errorType     ErrorType
	originalError error
}

func (et ErrorType) New(message string) error {
	return customError{errorType: et, originalError: errors.New(message)}
}

func (et ErrorType) Errorf(format string, args ...interface{}) error {
	return customError{errorType: et, originalError: errors.Errorf(format, args...)}
}

func (et ErrorType) Wrap(err error, message string) error {
	return customError{errorType: et, originalError: errors.Wrap(err, message)}
}

func (et ErrorType) Wrapf(err error, format string, args ...interface{}) error {
	return customError{errorType: et, originalError: errors.Wrapf(err, format, args...)}
}

func (e customError) Error() string {
	return e.originalError.Error()
}

func (e customError) Type() ErrorType {
	return e.errorType
}

func Wrap(err error, message string) error {
	we := errors.Wrap(err, message)
	if ce, ok := err.(typeGetter); ok {
		return customError{errorType: ce.Type(), originalError: we}
	}
	return customError{errorType: Unknown, originalError: we}
}

func Cause(err error) error {
	return errors.Cause(err)
}

func GetType(e error) ErrorType {
	for {
		if tg, ok := e.(typeGetter); ok {
			return tg.Type()
		}
		break
	}
	return Unknown
}

func GetHttpCode(e error) int {
	switch GetType(e) {
	case Unknown:
		return http.StatusBadRequest
	case Invalidated:
		return http.StatusBadRequest
	case Unauthorized:
		return http.StatusUnauthorized
	case Forbidden:
		return http.StatusForbidden
	case NotFound:
		return http.StatusNotFound
	case Conflict:
		return http.StatusConflict
	case Failed:
		return http.StatusInternalServerError
	default:
		return http.StatusBadRequest
	}
}
