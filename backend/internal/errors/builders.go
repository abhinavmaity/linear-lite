package errors

import "net/http"

func Validation(message string, fields FieldErrors) *AppError {
	return &AppError{
		Status:  http.StatusBadRequest,
		Code:    "validation_error",
		Message: message,
		Fields:  fields,
	}
}

func Unauthorized(message string) *AppError {
	return &AppError{
		Status:  http.StatusUnauthorized,
		Code:    "unauthorized",
		Message: message,
	}
}

func NotFound(message string) *AppError {
	return &AppError{
		Status:  http.StatusNotFound,
		Code:    "not_found",
		Message: message,
	}
}

func Conflict(message string, fields FieldErrors) *AppError {
	return &AppError{
		Status:  http.StatusConflict,
		Code:    "conflict",
		Message: message,
		Fields:  fields,
	}
}

func Internal(message string) *AppError {
	return &AppError{
		Status:  http.StatusInternalServerError,
		Code:    "internal_error",
		Message: message,
	}
}

func NotImplemented(message string) *AppError {
	return &AppError{
		Status:  http.StatusNotImplemented,
		Code:    "not_implemented",
		Message: message,
	}
}
