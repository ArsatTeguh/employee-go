package helper

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetUser(ctx *gin.Context) (JWTResponse, error) {
	user := ctx.MustGet("user").(JWTResponse)

	if user.Email == "" {
		response := &WithoutData{
			Code:    401,
			Message: "Sesion Berakhir",
		}
		response.Response(ctx)

		http.SetCookie(ctx.Writer, &http.Cookie{
			Name:     "token",
			Value:    "",
			Expires:  time.Unix(0, 0), // Set expires to the past
			HttpOnly: true,
			Path:     "/",
		})

		return user, gin.Error{Err: errors.New(user.Email)}
	}
	return user, nil
}
