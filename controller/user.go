package controller

import (
	"backend/dto"
	"backend/helper"
	"backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	DB *gorm.DB
}

type userServcie interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)
	GetOneUser(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
}

func NewUserService(db *gorm.DB) userServcie {
	return &User{
		DB: db,
	}
}

func (u *User) GetOneUser(ctx *gin.Context) {
	var user models.User
	users, err := helper.GetUser(ctx)

	if err != nil {
		return

	}

	if err := u.DB.First(&user, users.Id).Error; err != nil {
		helper.ErrorServer(err, ctx)
		return
	}

	response := &helper.WithData{
		Code:    200,
		Message: "success",
		Data:    user,
	}

	response.Response(ctx)
}

type regsiter struct {
	Id             *int64  `json:"id,omitempty"`
	Email          string  `json:"email" validate:"required,email"`
	Password       string  `validate:"required" json:"password"`
	RepeatPassword string  `validate:"required" json:"repeat_password"`
	Name           string  `validate:"required" json:"name"`
	Address        string  `validate:"required" json:"address"`
	Role           *string `json:"role,omitempty"`
}

func (u *User) Register(ctx *gin.Context) {

	var payload regsiter

	defaultUser := "karyawan"

	if err := dto.ValidationPayload(&payload, ctx); err != nil {
		return
	}

	if payload.Password != payload.RepeatPassword {
		response := &helper.WithoutData{
			Code:    400,
			Message: "Password and Repeat no Match",
		}
		response.Response(ctx)
		return
	}

	if payload.Role == nil {
		payload.Role = &defaultUser
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		response := &helper.WithoutData{
			Code:    400,
			Message: err.Error(),
		}
		response.Response(ctx)
		return
	}

	payload.Password = string(hashedPassword)

	user := models.User{
		Role:     payload.Role,
		Email:    payload.Email,
		Password: payload.Password,
	}

	if err := u.DB.Create(&user).Error; err != nil {
		response := &helper.WithoutData{
			Code:    500,
			Message: err.Error(),
		}
		response.Response(ctx)
		return
	}

	employee := models.Employee{
		Name:    payload.Name,
		Address: payload.Address,
		Status:  *user.Role,
		Email:   user.Email,
		UserId:  user.Id,
	}

	if err := u.DB.Create(&employee).Error; err != nil {
		response := &helper.WithoutData{
			Code:    500,
			Message: err.Error(),
		}
		response.Response(ctx)
		return
	}

	wallet := models.Wallet{
		EmployeeId: employee.Id,
	}

	if err := u.DB.Create(&wallet).Error; err != nil {
		response := &helper.WithoutData{
			Code:    500,
			Message: err.Error(),
		}
		response.Response(ctx)
		return
	}

	response := &helper.WithData{
		Code:    201,
		Message: "success",
		Data:    user,
	}
	response.Response(ctx)
}

type payload struct {
	Email    string `validate:"required,email" json:"email"`
	Password string `validate:"required" json:"password"`
}

func (u *User) Login(ctx *gin.Context) {
	var user models.User
	var p payload

	if err := dto.ValidationPayload(&p, ctx); err != nil {
		return
	}

	if err := u.DB.Where("email = ?", p.Email).First(&user).Error; err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "User email dan password tidak sesusai"})
		return
	}

	// math passoword
	errorHash := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(p.Password))

	if errorHash != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "User email dan password tidak sesusai"})
		return
	}

	token, err := helper.CreateToken(user.Email, user.Id, *user.Role)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Terjadi kesalahan di accesstoken"})
		return
	}

	refreshToken, err := helper.CreateRefreshToken(user.Email, user.Id, *user.Role)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Terjadi kesalahan di refreshtoken"})
		return
	}

	user.Tokenize = &refreshToken // update value refresh token

	if updateRefreshToken := u.DB.Updates(&user).Error; updateRefreshToken != nil {
		helper.ErrorServer(updateRefreshToken, ctx)
		return
	}

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
		Secure:   false, // Set to true when using HTTPS in production
		SameSite: http.SameSiteLaxMode,
	})

	response := &helper.WithData{
		Code:    200,
		Message: "Berhasil Login",
		Data:    token,
	}

	response.Response(ctx)
}

func (u *User) Logout(ctx *gin.Context) {
	var personalUser models.User
	refreshToken, err := ctx.Cookie("refreshToken")

	if err != nil {
		ctx.JSON(402, gin.H{"message": "token tidak ada dalam cookie"})
		return
	}

	if err := u.DB.First(&personalUser, "tokenize = ?", refreshToken).Error; err != nil {
		ctx.AbortWithError(403, err)
		return
	}

	ctx.SetCookie("refreshToken", "", -1, "", "localhost", false, true)

	personalUser.Tokenize = nil

	if updateRefreshToken := u.DB.Select("Tokenize").Updates(&personalUser).Error; updateRefreshToken != nil {
		helper.ErrorServer(updateRefreshToken, ctx)
		return
	}

	response := &helper.WithoutData{
		Code:    200,
		Message: "Berhasil Log out",
	}

	response.Response(ctx)
}

func (u *User) GetAll(ctx *gin.Context) {

	validate := helper.Premission(ctx)

	if validate != nil {
		return
	}

	var user []models.User

	if err := u.DB.Find(&user).Error; err != nil {
		helper.ErrorServer(err, ctx)
		return
	}

	response := &helper.WithData{
		Code:    200,
		Message: "success",
		Data:    user,
	}

	response.Response(ctx)
}

var SecreetRefreshToken = models.GetEnv("SECREETREFRESHTOKEN", "SECREET TIDAK DITEMUKAN")

func (u *User) RefreshToken(ctx *gin.Context) {
	var user models.User
	var customeclaims helper.MyCustomClaims

	refreshToken, err := ctx.Cookie("refreshToken")

	if err != nil {
		ctx.JSON(402, gin.H{"message": "token tidak ada dalam cookie"})
		return
	}

	if err := u.DB.First(&user, "tokenize = ?", refreshToken).Error; err != nil {
		ctx.JSON(402, gin.H{"message": "refresh token  tidak ditemukan"})
		return
	}

	token, err := jwt.ParseWithClaims(
		refreshToken,
		&customeclaims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SecreetRefreshToken), nil // Compare with the secret
		},
	)

	if err != nil {
		user.Tokenize = nil
		if updateRefreshToken := u.DB.Select("Tokenize").Updates(&user).Error; updateRefreshToken != nil {
			helper.ErrorServer(updateRefreshToken, ctx)
			return
		}
		ctx.JSON(403, gin.H{"message": "Refresh Token expired"})
		return
	}

	_, ok := token.Claims.(*helper.MyCustomClaims)

	if !ok || !token.Valid {
		ctx.JSON(403, gin.H{"message": "token no claims"})
		return
	}

	newToken, err := helper.CreateToken(user.Email, user.Id, *user.Role)

	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	response := &helper.WithData{
		Code:    200,
		Message: "success",
		Data:    newToken,
	}

	response.Response(ctx)
}
