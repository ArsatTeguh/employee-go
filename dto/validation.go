package dto

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ValidationPayload(payload interface{}, ctx *gin.Context) error {

	validate := validator.New()

	ctx.ShouldBindJSON(&payload)

	if err := validate.Struct(payload); err != nil {
		ctx.AbortWithStatusJSON(400, err.Error())
		return errors.New(err.Error())
	}

	return nil

}
