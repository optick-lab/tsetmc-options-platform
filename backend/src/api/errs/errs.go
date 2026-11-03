package errs

import (
	"errors"
	"fmt"
	"net/http"
)

// Sentinel errors
var (
	ErrNotFound       = errors.New("resource not found")
	ErrInternalServer = errors.New("internal server error")
	ErrBadRequest     = errors.New("bad request")
	ErrDatabaseDown   = errors.New("database connection unavailable")
)

// AppError represents an application-level structured error.
type AppError struct {
	Key       string `json:"key"`
	Status    int    `json:"status"`
	ErrorCode int    `json:"error_code"`
}

func (e AppError) Error() string {
	return fmt.Sprintf("[%d:%d] %s", e.Status, e.ErrorCode, e.Key)
}

// New creates a new AppError.
func New(key string, status, errorCode int) AppError {
	return AppError{
		Key:       key,
		Status:    status,
		ErrorCode: errorCode,
	}
}

// Common AppErrors
var (
	AppErrorNotFound = AppError{
		Key:       "not_found",
		Status:    http.StatusNotFound,
		ErrorCode: 40400,
	}
	AppErrorInternal = AppError{
		Key:       "internal_server_error",
		Status:    http.StatusInternalServerError,
		ErrorCode: 50000,
	}
)
