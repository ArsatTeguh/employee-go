package helper

import (
	"backend/models"
	"fmt"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var Secreet = models.GetEnv("SECREET", "SECREET TIDAK DITEMUKAN")

type MyCustomClaims struct {
	Id    int64  `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

type JWTResponse struct {
	Id    int64  `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	Exp   int64  `json:"exp"`
	Nbf   int64  `json:"nbf"`
	Iat   int64  `json:"iat"`
}

func CreateToken(email string, id int64, role string) (string, error) {
	claims := MyCustomClaims{
		id,
		email,
		role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(Secreet))
}

func ExtractToken(c *gin.Context) string {
	cookie, err := c.Request.Cookie("token")
	if err != nil {
		c.JSON(400, gin.H{"message": "token tidak ada dalam cookie"})
		return ""
	}

	return cookie.Value
}

func TokenValid(c *gin.Context) (JWTResponse, error) {
	tokenString := ExtractToken(c)
	j := JWTResponse{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&MyCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(Secreet), nil
		},
	)

	if err != nil {
		return j, fmt.Errorf("Unauthentication1")
	}

	claims, ok := token.Claims.(*MyCustomClaims)

	if !ok || !token.Valid {
		return j, fmt.Errorf("Unauthentication2")
	}
	claim := JWTResponse{
		Id:    claims.Id,
		Email: claims.Email,
		Role:  claims.Role,
		Exp:   claims.ExpiresAt.Unix(),
		Nbf:   claims.ExpiresAt.UnixMicro(),
		Iat:   claims.ExpiresAt.UnixMilli(),
	}
	return claim, nil
}
