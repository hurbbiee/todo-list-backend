package response

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/hurbbiee/todo-list-backend/internal/shared/enum"
)

type ErrorResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Detail  string            `json:"detail,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}

type AppError struct {
	HTTPCode int
	Code     int
	Message  string
	Detail   string
	Errors   map[string]string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(httpCode int, def enum.ErrorDef) *AppError {
	return &AppError{
		HTTPCode: httpCode,
		Code:     def.Code,
		Message:  def.Message,
	}
}

func ValidationError(def enum.ErrorDef, errs map[string]string) *AppError {
	return &AppError{
		HTTPCode: http.StatusBadRequest,
		Code:     def.Code,
		Message:  def.Message,
		Errors:   errs,
	}
}

func Unauthorized(def enum.ErrorDef) *AppError {
	return NewAppError(http.StatusUnauthorized, def)
}

func InternalError(def enum.ErrorDef) *AppError {
	return NewAppError(http.StatusInternalServerError, def)
}

func FromError(err error) (int, ErrorResponse) {
	if appErr, ok := err.(*AppError); ok {
		return appErr.HTTPCode, ErrorResponse{
			Code:    appErr.Code,
			Message: appErr.Message,
			Detail:  appErr.Detail,
			Errors:  appErr.Errors,
		}
	}

	if fiberErr, ok := err.(*fiber.Error); ok {
		return fiberErr.Code, ErrorResponse{
			Code:    fiberErr.Code,
			Message: fiberErr.Message,
		}
	}

	return http.StatusInternalServerError, ErrorResponse{
		Code:    enum.Internal.Code,
		Message: enum.Internal.Message,
	}
}

func BadRequest(def enum.ErrorDef) *AppError {
	return NewAppError(http.StatusBadRequest, def)
}
