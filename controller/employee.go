package controller

import (
	"backend/dto"
	"backend/helper"
	"backend/models"
	"fmt"
	"strconv"

	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type employeeController struct {
	DB *gorm.DB
}

type ServiceEmployee interface {
	GetAllEmployee(ctx *gin.Context)
	GetOneEmployee(ctx *gin.Context)
	SaveEmployee(ctx *gin.Context)
	UploadProfile(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

func NewServiceEmployee(db *gorm.DB) ServiceEmployee {
	return &employeeController{
		DB: db,
	}
}

func (e *employeeController) GetAllEmployee(ctx *gin.Context) {

	valid := helper.Premission(ctx)

	if valid != nil {
		return
	}

	var emp []models.Employee
	search := ctx.Query("name")

	var totalCount int64
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1")) // convert string ke int
	sizePage, _ := strconv.Atoi(ctx.DefaultQuery("sizePage", "5"))

	offset := (page - 1) * sizePage
	query := e.DB.Model(&emp).Preload("Wallet").Preload("Position").Preload("Position.Attedance")

	if search != "" {
		emailPattern := "%" + strings.ToLower(search) + "%"
		query = query.Where("lower(name) LIKE ?", emailPattern)
	}

	query.Count(&totalCount)
	query.Offset(offset).Limit(sizePage).Find(&emp)

	if len(emp) == 0 {
		response := &helper.WithoutData{
			Code:    400,
			Message: "Data kosong",
		}
		response.Response(ctx)
		return
	}

	response := &helper.WithData{
		Code:    200,
		Message: "Ok",
		Data: map[string]any{
			"employees":  emp,
			"total":      totalCount,
			"page":       page,
			"sizePage":   sizePage,
			"totalPages": (totalCount + int64(sizePage) - 1) / int64(sizePage),
		},
	}
	response.Response(ctx)

}

func (e *employeeController) GetOneEmployee(ctx *gin.Context) {
	user, err := helper.GetUser(ctx)

	if err != nil {
		return
	}

	employee, err := helper.EmployeeExist(user.Id, e.DB)

	if err != nil {
		response := &helper.WithoutData{
			Code:    500,
			Message: err.Error(),
		}
		response.Response(ctx)
		return
	}

	joiningDuration := helper.MonthsSinceJoined(employee.CreatedAt)

	response := &helper.WithData{
		Code:    200,
		Message: "Ok",
		Data: map[string]any{
			"employee":        employee,
			"joiningDuration": fmt.Sprintf("%d Month", joiningDuration),
		},
	}
	response.Response(ctx)

}

func (e *employeeController) SaveEmployee(ctx *gin.Context) {
	if valid := helper.Premission(ctx); valid != nil {
		return
	}

	var payload dto.CreatePayload
	if err := dto.ValidationPayload(&payload, ctx); err != nil {
		return
	}

	employee := payload.PayloadEmployee()

	if err := e.DB.Create(&employee).Error; err != nil {
		response := &helper.WithoutData{
			Code:    500,
			Message: err.Error(),
		}
		response.Response(ctx)
		return
	}

	wallet := payload.PayloadWallet(employee.Id)

	if err := e.DB.Create(&wallet).Error; err != nil {
		response := &helper.WithoutData{
			Code:    500,
			Message: "walet tidak dibuat",
		}

		response.Response(ctx)
		e.DB.Delete(&models.Employee{}, employee.Id)
		return
	}

	response := &helper.WithData{
		Code:    200,
		Message: "Data berhasil disimpan",
		Data:    employee,
	}
	response.Response(ctx)
}

func (e *employeeController) Update(ctx *gin.Context) {
	user, errors := helper.GetUser(ctx)

	if errors != nil {
		return
	}

	var body dto.CreatePayload
	var employee models.Employee

	if err := dto.ValidationPayload(&body, ctx); err != nil {
		return
	}

	_, err := helper.EmployeeExist(user.Id, e.DB)

	if err != nil {
		response := &helper.WithoutData{
			Code:    500,
			Message: err.Error(),
		}
		response.Response(ctx)
		return
	}

	query := e.DB.Model(&employee).Where("id = ?", user.Id) // update many values
	payload_em := body.PayloadEmployee()
	if err := query.Updates(&payload_em).Error; err != nil {
		response := &helper.WithoutData{
			Code:    400,
			Message: err.Error(),
		}
		response.Response(ctx)
		return
	}

	qw := e.DB.Model(&models.Wallet{}).Where("employee_id = ?", user.Id)
	payload_wl := body.PayloadWallet(user.Id)

	if err := qw.Updates(&payload_wl).Error; err != nil {
		response := &helper.WithoutData{
			Code:    500,
			Message: err.Error(),
		}

		response.Response(ctx)
		return
	}

	response := &helper.WithoutData{
		Code:    200,
		Message: "Update Success",
	}
	response.Response(ctx)
}

func (e *employeeController) Delete(ctx *gin.Context) {
	if valid := helper.Premission(ctx); valid != nil {
		return
	}

	var employee models.Employee
	var user models.User
	param := ctx.Param("id")
	id, _ := strconv.Atoi(param)

	if err := e.DB.Delete(&employee, id).Error; err != nil {
		helper.ErrorServer(err, ctx)
		return
	}

	if err := e.DB.Delete(&user, employee.Id).Error; err != nil {
		helper.ErrorServer(err, ctx)
		return
	}

	response := &helper.WithoutData{
		Code:    200,
		Message: "Delete Success",
	}

	response.Response(ctx)
}

func (e *employeeController) UploadProfile(ctx *gin.Context) {
	var employee models.Employee

	user, err := helper.GetUser(ctx)

	file, error := ctx.FormFile("image")

	if error != nil {
		response := &helper.WithoutData{
			Code:    400,
			Message: error.Error(),
		}

		response.Response(ctx)
		return
	}

	if err != nil {
		response := &helper.WithoutData{
			Code:    500,
			Message: err.Error(),
		}

		response.Response(ctx)
		return
	}

	urlStatic := "assets/" + file.Filename

	if err := ctx.SaveUploadedFile(file, urlStatic); err != nil {
		response := &helper.WithoutData{
			Code:    500,
			Message: err.Error(),
		}

		response.Response(ctx)
		return
	}

	if err := e.DB.Model(&employee).Where("id = ?", user.Id).Update("picture", "/"+urlStatic).Error; err != nil {
		response := &helper.WithoutData{
			Code:    400,
			Message: err.Error(),
		}

		response.Response(ctx)
		return
	}

	response := &helper.WithoutData{
		Code:    200,
		Message: "Upload Success",
	}

	response.Response(ctx)
}
