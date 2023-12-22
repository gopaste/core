package typesystem

import (
	"fmt"
	"net/http"
)

var (
	BadRequest  = NewHttpError("Bad request occurred", "Missing required parameters", http.StatusBadRequest)
	NotFound    = NewHttpError("Resource not found", "The requested resource does not exist", http.StatusNotFound)
	ServerError = NewHttpError(
		"Internal server error",
		"An unexpected error occurred on the server",
		http.StatusInternalServerError,
	)
	UserConflictError = NewHttpError("User conflict", "A user with the same email already exists", http.StatusConflict)

	Unauthorized      = NewHttpError("Unauthorized", "Authentication failed or missing credentials", http.StatusUnauthorized)
	TokenInvalidError = NewHttpError("Token invalid", "The provided token is not valid", http.StatusUnauthorized)
	TokenExpiredError = NewHttpError("Token expired", "The provided token has expired", http.StatusUnauthorized)
)

type Http struct {
	Description string `json:"description,omitempty"`
	Metadata    string `json:"metadata,omitempty"`
	StatusCode  int    `json:"statusCode"`
}

func (e Http) Error() string {
	return fmt.Sprintf("description: %s,  metadata: %s", e.Description, e.Metadata)
}

func NewHttpError(description, metadata string, statusCode int) Http {
	return Http{
		Description: description,
		Metadata:    metadata,
		StatusCode:  statusCode,
	}
}
