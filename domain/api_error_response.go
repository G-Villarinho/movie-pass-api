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

var CannotBindPayloadAPIError = ErrorResponse{
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

var InternalServerErrorAPIError = ErrorResponse{
	StatusCode: http.StatusInternalServerError,
	Title:      "Internal Server Error",
	Details:    "Something went wrong on our end. Please try again later or contact support if the issue persists.",
	Errors:     nil,
}

func NewValidationErrorResponse(ctx echo.Context, statusCode int, validationErrors ValidationErrors) error {
	errorResponse := ErrorResponse{
		StatusCode: statusCode,
		Title:      "Validation Error",
		Details:    "One or more fields are invalid.",
		Errors:     convertToValidationErrorList(validationErrors),
	}

	return ctx.JSON(statusCode, errorResponse)
}

func NewCustomValidationErrorResponse(ctx echo.Context, statusCode int, validationErrors ValidationErrors, title, details string) error {
	errorResponse := ErrorResponse{
		StatusCode: statusCode,
		Title:      title,
		Details:    details,
		Errors:     convertToValidationErrorList(validationErrors),
	}

	return ctx.JSON(statusCode, errorResponse)
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

func CannotBindPayloadError(ctx echo.Context) error {
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

func InternalServerError(ctx echo.Context) error {
	errorResponse := ErrorResponse{
		StatusCode: http.StatusInternalServerError,
		Title:      "Internal Server Error",
		Details:    "Something went wrong on our end. Please try again later or contact support if the issue persists.",
		Errors:     nil,
	}
	return ctx.JSON(http.StatusUnprocessableEntity, errorResponse)
}
