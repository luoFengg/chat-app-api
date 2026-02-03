package web

// ApiResponse is a Wrapper for All API Response
type ApiResponse struct {
	Success bool `json:"success"`
	Message string `json:"message"`
	Data any `json:"data,omitempty"`
}	

// PaginationMeta for Response with Pagination
type PaginationMeta struct {
	CurrentPage int `json:"current_page"`
	PerPage int `json:"per_page"`
	TotalItems int64 `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

// PaginatedResponse for Response with Data + Pagination
type PaginatedResponse struct {
	Success bool `json:"success"`
	Message string `json:"message"`
	Data any `json:"data,omitempty"`
	Pagination PaginationMeta `json:"pagination"`
}

// CursorMeta for Response with Cursor Pagination
type CursorMeta struct {
	HasMore bool `json:"has_more"`
	NextCursor *string `json:"next_cursor,omitempty"`
	PrevCursor *string `json:"prev_cursor,omitempty"`	
}

// CursorResponse fir Response with Cursor Pagination
type CursorResponse struct {
	Success bool `json:"success"`
	Message string `json:"message"`
	Data any `json:"data,omitempty"`
	Cursor CursorMeta `json:"cursor"`
}

// ErrorResponse for Error Response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}
