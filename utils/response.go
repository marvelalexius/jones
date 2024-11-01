package utils

import (
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type (
	ValidationErrorMsg struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	}

	ErrorRes struct {
		Message string `json:"message"`
		Debug   error  `json:"debug,omitempty"`
		Errors  any    `json:"errors"`
	}

	SuccessRes struct {
		Message string `json:"message"`
		Data    any    `json:"data,omitempty"`
		Meta    any    `json:"meta,omitempty"`
	}
)

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Should be a valid email address"
	case "file":
		return "Should be a valid file"
	case "lte":
		return "Should be less than " + fe.Param()
	case "lt":
		return "Should be less than " + fe.Param() + "digits"
	case "gte":
		return "Should be greater than " + fe.Param()
	case "eqfield":
		return "Should be equal to " + fe.Param()
	case "contains":
		return "Should contains " + fe.Param()
	case "startsnotwith":
		return "Should not starts with " + fe.Param()
	case "len":
		return "Should have length " + fe.Param()
	case "oneof":
		return "Should be one of " + fe.Param()
	case "min":
		return "Should be greater than " + fe.Param()
	case "max":
		return "Should be less than " + fe.Param()
	case "number":
		return "Should be a number"
	case "numeric":
		return "Should be a number"
	case "boolean":
		return "Should be a boolean"
	case "gt":
		return "Should be greater than " + fe.Param() + " digits"
	default:
		return "Unknown error"
	}
}

func ValidationResponse(err error) []ValidationErrorMsg {
	var ve validator.ValidationErrors
	var marshalErr *json.UnmarshalTypeError

	// Check if the error is json unmarshal error
	if errors.As(err, &marshalErr) {
		return []ValidationErrorMsg{{marshalErr.Field, "Should be a valid " + marshalErr.Type.String()}}
	}

	if errors.As(err, &ve) {
		out := make([]ValidationErrorMsg, len(ve))
		for i, fe := range ve {
			out[i] = ValidationErrorMsg{fe.Field(), getErrorMsg(fe)}
		}
		return out
	}

	return nil
}

func ErrorResponse(c *gin.Context, code int, res ErrorRes) {
	c.JSON(code, res)
}

func SuccessResponse(c *gin.Context, code int, res SuccessRes) {
	c.JSON(code, res)
}
