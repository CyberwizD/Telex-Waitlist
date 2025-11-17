package response

import "github.com/gin-gonic/gin"

// ErrorResponse represents a standardized error payload.
type ErrorResponse struct {
	Error string `json:"error"`
}

// DataResponse wraps a single data payload.
type DataResponse[T any] struct {
	Data T `json:"data"`
}

// PaginatedResponse wraps a list payload with pagination metadata.
type PaginatedResponse[T any] struct {
	Data   []T   `json:"data"`
	Total  int64 `json:"total"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}

// JSONError writes a JSON error response with a status code.
func JSONError(c *gin.Context, status int, message string) {
	c.JSON(status, ErrorResponse{Error: message})
}

// JSONData writes a JSON response containing a single object.
func JSONData[T any](c *gin.Context, status int, data T) {
	c.JSON(status, DataResponse[T]{Data: data})
}

// JSONPage writes a JSON response containing a list and pagination metadata.
func JSONPage[T any](c *gin.Context, status int, data []T, total int64, limit, offset int) {
	c.JSON(status, PaginatedResponse[T]{
		Data:   data,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}
