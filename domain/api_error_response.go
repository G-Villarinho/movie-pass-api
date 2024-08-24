package domain

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	StatusCode int               `json:"status"`
	Title      string            `json:"title"`
	Details    string            `json:"details"`
	Errors     []ValidationError `json:"errors,omitempty"`
}

func NewValidationAPIErrorResponse(ctx echo.Context, statusCode int, validationErrors ValidationErrors) error {
	errorResponse := ErrorResponse{
		StatusCode: statusCode,
		Title:      "Validation Error",
		Details:    "One or more fields are invalid.",
		Errors:     convertToValidationErrorList(validationErrors),
	}

	return ctx.JSON(statusCode, errorResponse)
}

func NewCustomValidationAPIErrorResponse(ctx echo.Context, statusCode int, validationErrors ValidationErrors, title, details string) error {
	errorResponse := ErrorResponse{
		StatusCode: statusCode,
		Title:      title,
		Details:    details,
		Errors:     convertToValidationErrorList(validationErrors),
	}

	return ctx.JSON(statusCode, errorResponse)
}

func CannotBindPayloadAPIErrorResponse(ctx echo.Context) error {
	errorResponse := ErrorResponse{
		StatusCode: http.StatusUnprocessableEntity,
		Title:      "Unable to Process Request",
		Details:    "We encountered an issue while trying to process your request. The data you provided is not in the expected format.",
		Errors: []ValidationError{
			{
				Field:   "payload",
				Message: "The information provided is not correctly formatted or is missing required fields. Please review and try again.",
			},
		},
	}
	return ctx.JSON(http.StatusUnprocessableEntity, errorResponse)
}

func InternalServerAPIErrorResponse(ctx echo.Context) error {
	errorResponse := ErrorResponse{
		StatusCode: http.StatusInternalServerError,
		Title:      "Internal Server Error",
		Details:    "Something went wrong on our end. Please try again later or contact support if the issue persists.",
		Errors:     nil,
	}
	return ctx.JSON(http.StatusInternalServerError, errorResponse)
}

func AccessDeniedAPIErrorResponse(ctx echo.Context) error {
	errorResponse := ErrorResponse{
		StatusCode: http.StatusUnauthorized,
		Title:      "Access Denied",
		Details:    "You need to be logged in to access this resource.",
		Errors:     nil,
	}
	return ctx.JSON(http.StatusUnauthorized, errorResponse)
}

func convertToValidationErrorList(validationErrors ValidationErrors) []ValidationError {
	errorList := make([]ValidationError, 0, len(validationErrors))
	for field, message := range validationErrors {
		errorList = append(errorList, ValidationError{
			Field:   field,
			Message: message,
		})
	}
	return errorList
}
