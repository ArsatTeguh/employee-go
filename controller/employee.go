package controller

import (
	"backend/dto"
	"backend/helper"
	"backend/models"
	"fmt"
	"os"
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
	GetProfile(ctx *gin.Context)
	PopupEmployee(ctx *gin.Context)
}

func NewServiceEmployee(db *gorm.DB) ServiceEmployee {
	return &employeeController{
		DB: db,
	}
}

type popup struct {
	Name   string `json:"name" `
	Id     int64  `json:"id" `
	Status string `json:"status" `
}

func (e *employeeController) PopupEmployee(ctx *gin.Context) {
	valid := helper.Premission(ctx)

	if valid != nil {
		return
	}
	var employee models.Employee

	result := []popup{}

	if err := e.DB.Model(&employee).Select("id", "name", "status").Find(&result).Error; err != nil {
		helper.ErrorServer(err, ctx)
		return
	}

	response := &helper.WithData{
		Code:    200,
		Message: " Success",
		Data:    result,
	}
	response.Response(ctx)

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
	query := e.DB.Model(&emp).Preload("Wallet").Preload("Position").Preload("Project").Preload("Project.Attedance")

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
			"total":      len(emp),
			"totalAll":   totalCount,
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

	var body dto.UpdateEmployeeStruct
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

	query := e.DB.Model(&employee).Where("id = ?", user.Id).First(&employee) // update many values
	payload_em := body.UpdateEmployee()
	if err := query.Updates(&payload_em).Error; err != nil {
		response := &helper.WithoutData{
			Code:    400,
			Message: err.Error(),
		}
		response.Response(ctx)
		return
	}

	response := &helper.WithData{
		Code:    200,
		Message: "Update Success",
		Data:    employee,
	}
	response.Response(ctx)
}

func (e *employeeController) Delete(ctx *gin.Context) {
	// Check permissions
	if valid := helper.Premission(ctx); valid != nil {
		return
	}

	param := ctx.Param("id")

	if param == "" {
		ctx.AbortWithStatusJSON(400, map[string]string{"message": "Param Required"})
		return
	}

	id, _ := strconv.ParseInt(param, 10, 64)

	var user models.User

	if err := helper.AssociationEmploee(id, ctx, e.DB); err != nil {
		ctx.AbortWithStatusJSON(400, map[string]string{"message": "Failed delete association", "error": err.Error()})
		return
	}

	if err := e.DB.Delete(&user, id).Error; err != nil {
		helper.ErrorServer(err, ctx)
		return
	}

	// Respond with success
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

	query := e.DB.Model(&employee).Where("id = ?", user.Id).First(&employee)

	if employee.Picture != "" {
		// Implement file deletion logic here
		os.Remove(employee.Picture)
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

	if err := query.Update("picture", urlStatic).Error; err != nil {
		response := &helper.WithoutData{
			Code:    400,
			Message: err.Error(),
		}
		response.Response(ctx)
		return
	}

	response := &helper.WithData{
		Code:    200,
		Message: "Upload Success",
		Data:    employee.Picture,
	}

	response.Response(ctx)
}

func (e *employeeController) GetProfile(ctx *gin.Context) {
	id := ctx.Param("id")
	var employee models.Employee

	if err := e.DB.Model(&employee).Where("id = ?", id).First(&employee).Error; err != nil {
		helper.ErrorServer(err, ctx)
		return
	}

	response := &helper.WithData{
		Code:    200,
		Message: "Success",
		Data:    employee.Picture,
	}

	response.Response(ctx)
}
