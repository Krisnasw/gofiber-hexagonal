package helpers

import (
	"net/http"
	"time"
)

// APIResponse represents the standard API response structure
type APIResponse struct {
	Error      bool        `json:"error"`
	StatusCode int         `json:"code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Metadata   Metadata    `json:"metadata,omitempty"`
}

// Metadata represents response metadata
type Metadata struct {
	Timestamp int64     `json:"timestamp,omitempty"`
	Page      *PageInfo `json:"page,omitempty"`
	TraceID   string    `json:"trace_id,omitempty"`
	Version   string    `json:"version,omitempty"`
}

// PageInfo represents pagination information
type PageInfo struct {
	CurrentPage    int    `json:"current_page,omitempty"`
	PageSize       int    `json:"page_size,omitempty"`
	TotalRecords   int    `json:"total_records,omitempty"`
	TotalPages     int    `json:"total_pages,omitempty"`
	HasNext        bool   `json:"has_next,omitempty"`
	HasPrevious    bool   `json:"has_previous,omitempty"`
	NextCursor     string `json:"next_cursor,omitempty"`
	PreviousCursor string `json:"previous_cursor,omitempty"`
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
		Metadata: Metadata{
			Timestamp: getCurrentTimestamp(),
		},
	}
}

// SuccessResponseWithMetadata creates a standardized success response with metadata.
func SuccessResponseWithMetadata(data interface{}, code int, message string, metadata Metadata) APIResponse {
	if message == "" {
		message = getStatusMessage(code)
	}
	// Ensure timestamp is always set
	if metadata.Timestamp == 0 {
		metadata.Timestamp = getCurrentTimestamp()
	}
	return APIResponse{
		Error:      false,
		StatusCode: code,
		Message:    message,
		Data:       data,
		Metadata:   metadata,
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
		Metadata: Metadata{
			Timestamp: getCurrentTimestamp(),
		},
	}
}

// DetailedErrorResponse creates an error response with additional details
func DetailedErrorResponse(code int, message, details string) APIResponse {
	return APIResponse{
		Error:      true,
		StatusCode: code,
		Message:    message,
		Data:       ErrorResponseDetail{Code: code, Message: message, Details: details},
		Metadata: Metadata{
			Timestamp: getCurrentTimestamp(),
		},
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

// getCurrentTimestamp returns the current Unix timestamp
func getCurrentTimestamp() int64 {
	return time.Now().Unix()
}
