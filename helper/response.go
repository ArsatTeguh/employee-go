package helper

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WithData struct {
	Code    int
	Message string
	Data    any
}

func (d *WithData) Response(ctx *gin.Context) {
	ctx.JSON(
		d.Code,
		gin.H{
			"message": d.Message,
			"data":    d.Data,
		},
	)
}

type WithoutData struct {
	Code    int
	Message string
}

func (d *WithoutData) Response(ctx *gin.Context) {
	ctx.JSON(
		d.Code,
		gin.H{
			"message": d.Message,
		},
	)

}

func ErrorServer(e error, ctx *gin.Context) error {
	if e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(400, map[string]string{"message": "Data Not Found"})
			return errors.New(e.Error())
		}
		ctx.AbortWithStatusJSON(500, map[string]string{"meesage": e.Error()})
		return errors.New(e.Error())
	}
	return nil
}
