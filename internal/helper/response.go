package helpers

import (
	"net/http"
)

// APIResponse represents the standard API response structure
type APIResponse struct {
	Error      bool        `json:"error"`
	StatusCode int         `json:"code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

// ErrorResponseDetail provides detailed error information
type ErrorResponseDetail struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// SuccessResponse creates a standardized success response.
func SuccessResponse(data interface{}, code int, message string) APIResponse {
	if message == "" {
		message = getStatusMessage(code)
	}
	return APIResponse{
		Error:      false,
		StatusCode: code,
		Message:    message,
		Data:       data,
	}
}

// ErrorResponse creates a standardized error response.
func ErrorResponse(data interface{}, code int, message string) APIResponse {
	if message == "" {
		message = getStatusMessage(code)
	}
	return APIResponse{
		Error:      true,
		StatusCode: code,
		Message:    message,
		Data:       data,
	}
}

// DetailedErrorResponse creates an error response with additional details
func DetailedErrorResponse(code int, message, details string) APIResponse {
	return APIResponse{
		Error:      true,
		StatusCode: code,
		Message:    message,
		Data:       ErrorResponseDetail{Code: code, Message: message, Details: details},
	}
}

// getStatusMessage maps HTTP status codes to default messages.
func getStatusMessage(code int) string {
	switch code {
	case http.StatusOK:
		return "OK"
	case http.StatusCreated:
		return "Created"
	case http.StatusNoContent:
		return "No Content"
	case http.StatusBadRequest:
		return "Bad Request"
	case http.StatusUnauthorized:
		return "Unauthorized"
	case http.StatusForbidden:
		return "Forbidden"
	case http.StatusNotFound:
		return "Not Found"
	case http.StatusUnprocessableEntity:
		return "Unprocessable Entity"
	case http.StatusInternalServerError:
		return "Internal Server Error"
	default:
		return "Unknown Code"
	}
}
