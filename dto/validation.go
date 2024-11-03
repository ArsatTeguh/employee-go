package dto

import (
	"errors"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ValidationPayload(payload interface{}, ctx *gin.Context) error {
	validate := validator.New()

	// Bind the incoming JSON to the payload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{"message": "Invalid JSON"})
		return errors.New(err.Error())
	}

	// Check if the payload is a slice
	v := reflect.ValueOf(payload).Elem()

	if v.Kind() == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			// Validate each item in the slice
			item := v.Index(i)
			if err := validate.Struct(item); err != nil {
				ctx.AbortWithStatusJSON(400, gin.H{"message": "Not Validate Character", "error": err.Error()})
				return errors.New(err.Error())
			}
		}
	} else {
		// Validate single object
		if err := validate.Struct(payload); err != nil {
			ctx.AbortWithStatusJSON(400, gin.H{"message": "Not Validate Character", "error": err.Error()})
			return errors.New(err.Error())
		}
	}

	return nil
}
