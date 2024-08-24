package domain

import "github.com/labstack/echo/v4"

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Title   string            `json:"title"`
	Details string            `json:"details"`
	Errors  []ValidationError `json:"errors,omitempty"`
}

func NewValidationErrorResponse(validationErrors ValidationErrors) ErrorResponse {
	errorList := make([]ValidationError, 0, len(validationErrors))

	for field, message := range validationErrors {
		errorList = append(errorList, ValidationError{
			Field:   field,
			Message: message,
		})
	}

	return ErrorResponse{
		Title:   "Validation Error",
		Details: "One or more fields are invalid.",
		Errors:  errorList,
	}
}

func WriteAPIErrorResponse(ctx echo.Context, statusCode int, errorResponse ErrorResponse) error {
	return ctx.JSON(statusCode, errorResponse)
}
