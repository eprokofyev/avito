package errors

import (
	"avito/internal/pkg/response"
	"fmt"
	"github.com/pkg/errors"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

const(
	NoType = ErrorType(iota)
	BalanceNotFound
	BadUserData
	BadData
	InsufficientFunds
)

type ErrorType uint

type Error struct {
	errorType ErrorType
	originalError error
	message string
}

func (error Error) Error() string {
	return error.originalError.Error()
}

func (e ErrorType) New(msg string) error {
	return Error{errorType: e,
		originalError: New(msg),
	}
}

func (e ErrorType) Newf(msg string, args ...interface{}) error {
	err := fmt.Errorf(msg, args...)
	return Error{errorType: e, originalError: err}
}

func (e ErrorType) Wrap(err error, msg string) error {
	return e.Wrapf(err, msg)
}

func (e ErrorType) Wrapf(err error, msg string, args ...interface{}) error {
	newErr := errors.Wrapf(err, msg, args...)

	return Error{errorType: e, originalError: newErr}
}

func New(msg string) error {
	return Error{errorType: NoType, originalError: errors.New(msg)}
}

func Newf(msg string, args ...interface{}) error {
	return Error{errorType: NoType, originalError: New(fmt.Sprintf(msg, args...))}
}

func Wrap(err error, msg string) error {
	return Wrapf(err, msg)
}


func Wrapf(err error, msg string, args ...interface{}) error {
	wrappedError := errors.Wrapf(err, msg, args...)
	if customErr, ok := err.(Error); ok {
		return Error{
			errorType: customErr.errorType,
			originalError: wrappedError,
			message: customErr.message,
		}
	}

	return Error{errorType: NoType, originalError: wrappedError}
}

func GetType(err error) ErrorType {
	if customErr, ok := err.(Error); ok {
		return customErr.errorType
	}

	return NoType
}

func WithMessage(err error, message string) error {
	if customErr, ok := err.(Error); ok {
		customErr.message = message
		return customErr
	}

	return Error{errorType: NoType, originalError: err, message: message}
}

func ErrorHandler(ctx *routing.Context, err error) error {
	var status int
	var message string

	switch GetType(err) {
	case InsufficientFunds:
		status = fasthttp.StatusConflict
		message = "Insufficient funds to write off"
	case BadUserData:
		status = fasthttp.StatusNotAcceptable
		message = "sender_id and recipient_id can't be zero at the same request"
	case BadData:
		status = fasthttp.StatusBadRequest
		message = "Request data is not correct"
	case BalanceNotFound:
		status = fasthttp.StatusNotFound
		message = "User's balance is not found"
	default:
		status = fasthttp.StatusInternalServerError
		message = "Internal server error"
	}

	return response.Respond(ctx, status, map[string]interface{} {"message":message})
}
