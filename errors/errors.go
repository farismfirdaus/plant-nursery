package errors

import (
	"net/http"
)

// Error contain application error message
type Error struct {
	// Message hold error message
	Message string

	// Code hold error code.
	// each error should have unique code.
	Code int

	// StatusCode hold http status code.
	StatusCode int
}

// Error return error message
func (e Error) Error() string {
	return e.Message
}

var (
	// BadRequest client request is not valid
	BadRequest = New("bad request", 100, http.StatusBadRequest)

	// Unauthorized client request is not authorized
	Unauthorized = New("unauthorized", 101, http.StatusUnauthorized)

	// InvalidCustomerFullName client request empty customer full name
	InvalidCustomerFullName = New("full name should not be empty", 500, http.StatusBadRequest)

	// InvalidEmail client request invalid email
	InvalidEmail = New("invalid email", 501, http.StatusBadRequest)

	// InvalidPassword client request empty password
	InvalidPassword = New("password should not be empty", 502, http.StatusBadRequest)

	// InvalidEmailAlreadyTaken client request email that has already been taken
	InvalidEmailAlreadyTaken = New("email already been taken", 503, http.StatusBadRequest)

	// InvalidEmailNotFound client request email not exists
	InvalidEmailNotFound = New("email not found", 504, http.StatusNotFound)

	// InvalidPasswordNotMatch client request email not exists
	InvalidPasswordNotMatch = New("password not match", 505, http.StatusBadRequest)

	// InvalidCustomerID client request invalid customer id
	InvalidCustomerID = New("invalid customer id", 506, http.StatusBadRequest)

	// InvalidPlantID client request invalid plant id
	InvalidPlantID = New("invalid plant id", 507, http.StatusBadRequest)

	// InvalidStockNotAvailable client request stock not available
	InvalidStockNotAvailable = New("stock not available", 508, http.StatusBadRequest)

	// InvalidCartNotFound client request cart not found
	InvalidCartNotFound = New("cart not found", 509, http.StatusNotFound)
)

// New create a new error
func New(message string, code int, statusCode int) error {
	return Error{
		Message:    message,
		Code:       code,
		StatusCode: statusCode,
	}
}
