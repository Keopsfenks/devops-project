package dto

// Result represents a standardized response structure for application operations
type Result[T any] struct {
	Data         *T     `json:"data,omitempty"`
	IsSuccessful bool   `json:"is_successful"`
	ErrorMessage string `json:"error_message,omitempty"`
	StatusCode   int    `json:"status_code"`
}

func NewResult[T any](data *T, isSuccessful bool, errorMessage string, statusCode int) *Result[T] {
	return &Result[T]{
		Data:         data,
		IsSuccessful: isSuccessful,
		ErrorMessage: errorMessage,
		StatusCode:   statusCode,
	}
}
