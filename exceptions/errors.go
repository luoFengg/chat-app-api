package exceptions

// NotFoundError for Resource Not Found
type NotFoundError struct {
	Message string
}

func NewNotFoundError(message string) NotFoundError {
	return NotFoundError{Message: message}
}

func (e NotFoundError) Error() string {
	return e.Message
}

// BadRequestError for Invalid Request
type BadRequestError struct {
	Message string
}

func NewBadRequestError(message string) BadRequestError {
	return BadRequestError{Message: message}
}

func (e BadRequestError) Error() string {
	return e.Message
}

// UnauthorizedError for Unauthorized Request
type UnauthorizedError struct {
	Message string
}

func NewUnauthorizedError(message string) UnauthorizedError {
	return UnauthorizedError{Message: message}
}

func (e UnauthorizedError) Error() string {
	return e.Message
}

// ConflictError for Data Conflict (duplicate data)
type ConflictError struct {
	Message string
}

func NewConflictError(message string) ConflictError {
	return ConflictError{Message: message}
}

func (e ConflictError) Error() string {
	return e.Message
}

// ForbiddenError for the Forbidden Access
type ForbiddenError struct {
	Message string
}

func NewForbiddenError(message string) ForbiddenError {
	return ForbiddenError{Message: message}
}

func (e ForbiddenError) Error() string {
	return e.Message
}

