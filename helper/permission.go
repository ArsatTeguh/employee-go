package helper

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Premission(ctx *gin.Context) error {
	user := ctx.MustGet("user").(JWTResponse)

	if user.Role == "karyawan" || user.Role == "" {
		response := &WithoutData{
			Code:    400,
			Message: "Akses tidak dizinkan",
		}
		response.Response(ctx)
		return fmt.Errorf("")
	}

	return nil

}
