package helper

import (
	"reflect"
)

var Environment = "development"

// Response represents a standard API response
type Response struct {
	HttpCode   int      `json:"httpCode,omitempty"`
	ErrorCode  int      `json:"errorCode,omitempty"`
	ErrorName  string   `json:"errorName,omitempty"`
	Message    string   `json:"message,omitempty"`
	Data       any      `json:"data,omitempty"`
	StackTrace []string `json:"stack,omitempty"`
	Details    string   `json:"details,omitempty"`
}

// Pagination represents pagination parameters
type Pagination struct {
	Page       int `json:"page" form:"page"`
	PageSize   int `json:"page_size" form:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// Filter represents query filter parameters
type Filter struct {
	Field    string `json:"field" form:"field"`
	Operator string `json:"operator" form:"operator"`
	Value    any    `json:"value" form:"value"`
}

// Sort represents sorting parameters
type Sort struct {
	Field     string `json:"field" form:"field"`
	Direction string `json:"direction" form:"direction"` // "asc" or "desc"
}

// NewSuccessResponse creates a success response
func NewSuccessResponse(data any) Response {
	return Response{
		Data: data,
	}
}

// Error implements the error interface
func (e *Response) Error() string {
	return e.Message
}

func WebResponse(response *Response) *Response {
	if response.ErrorCode > 0 && response.ErrorName != "" && Environment == "development" {
		// response.StackTrace = strings.Split(string(debug.Stack()), "\n")
	}

	return response
}

// NewPaginatedResponse creates a paginated response
func NewPaginatedResponse(data any, pagination Pagination) Response {
	return Response{
		Data: map[string]any{
			"items":      data,
			"pagination": pagination,
		},
	}
}

// Paginate applies pagination to a slice
func Paginate(data any, page, pageSize int) (any, Pagination) {
	s := reflect.ValueOf(data)
	if s.Kind() != reflect.Slice && s.Kind() != reflect.Array {
		return data, Pagination{}
	}

	total := s.Len()
	if total == 0 {
		return []any{}, Pagination{}
	}

	// Calculate pagination
	totalPages := (total + pageSize - 1) / pageSize

	// Validate page and pageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	// Calculate start and end indices
	start := (page - 1) * pageSize
	if start > total {
		start = total
	}

	end := start + pageSize
	if end > total {
		end = total
	}

	// Get paginated data
	result := reflect.MakeSlice(s.Type(), end-start, end-start)
	for i := start; i < end; i++ {
		result.Index(i - start).Set(s.Index(i))
	}

	return result.Interface(), Pagination{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
}
