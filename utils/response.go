package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func DataResponse(response interface{}) gin.H {
	return gin.H{"data": response}
}

func ErrorResponse(err error) gin.H {
	return gin.H{"errors": err.Error()}
}

func ValidationErrorsResponse(err error) gin.H {
	var errors []map[string]string
	for _, err := range err.(validator.ValidationErrors) {
		field := ToSnakeCase(err.Field())
		errorField := map[string]string{
			"field":   field,
			"message": fmt.Sprintf("%s is invalid with %s tag", field, err.Tag()),
		}
		errors = append(errors, errorField)
	}

	return gin.H{"errors": errors}
}
