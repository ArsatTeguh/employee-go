package controller

import (
	"backend/dto"
	"backend/helper"
	"backend/models"
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LeaveController struct {
	DB *gorm.DB
}

func (l *LeaveController) GetAll(ctx *gin.Context) {
	err := helper.Premission(ctx)
	if err != nil {
		return
	}

	var leave []models.Leave
	search := ctx.Query("status")

	var totalCount int64
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1")) // convert string ke int
	sizePage, _ := strconv.Atoi(ctx.DefaultQuery("sizePage", "5"))

	offset := (page - 1) * sizePage
	query := l.DB.Model(&leave).Preload("Employee", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	})

	if search != "" {
		emailPattern := "%" + strings.ToLower(search) + "%"
		query = query.Where("status LIKE ?", emailPattern)
	}

	query.Count(&totalCount)
	query.Offset(offset).Limit(sizePage).Find(&leave)

	if len(leave) == 0 {
		response := &helper.WithoutData{
			Code:    400,
			Message: "Data kosong",
		}
		response.Response(ctx)
		return
	}

	response := &helper.WithData{
		Code:    200,
		Message: "Success Get All Leaving",
		Data: map[string]any{
			"leave":      leave,
			"total":      len(leave),
			"totalAll":   totalCount,
			"page":       page,
			"sizePage":   sizePage,
			"totalPages": (totalCount + int64(sizePage) - 1) / int64(sizePage),
		},
	}
	response.Response(ctx)
}

func (l *LeaveController) Created(ctx *gin.Context) {
	var body dto.RequestLeave

	id := ctx.Param("id") // id of employee in request leave
	no, _ := strconv.ParseInt(id, 10, 64)

	employee, err := helper.EmployeeExist(no, l.DB)

	if err != nil {
		return
	}

	if err := dto.ValidationPayload(&body, ctx); err != nil {
		return
	}

	res := body.SavePosition(employee.Id)

	if err := l.DB.Create(&res).Error; err != nil {
		response := &helper.WithoutData{
			Code:    400,
			Message: err.Error(),
		}
		response.Response(ctx)
		return
	}

	go helper.SendEmail(employee)

	response := &helper.WithData{
		Code:    201,
		Message: "Insert",
		Data:    res,
	}
	response.Response(ctx)
}

type BodyJson struct {
	Status     string `json:"status" binding:"required"`
	EmployeeId int64  `json:"employee_id" binding:"required"`
}

func (l *LeaveController) Approve(ctx *gin.Context) {
	err := helper.Premission(ctx)
	if err != nil {
		return
	}

	var body BodyJson
	var leave models.Leave
	param := ctx.Param("id")
	id, _ := strconv.Atoi(param)
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response := &helper.WithoutData{
			Code:    400,
			Message: err.Error(),
		}
		response.Response(ctx)
		return
	}

	if err := l.DB.Where("id = ? & employee_id = ?", id, body.EmployeeId).First(&leave).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := &helper.WithoutData{
				Code:    400,
				Message: err.Error(),
			}
			response.Response(ctx)
			return
		}
		response := &helper.WithoutData{
			Code:    500,
			Message: err.Error(),
		}
		response.Response(ctx)
		return
	}

	leave.Status = &body.Status

	if err := l.DB.Updates(&leave).Error; err != nil {
		response := &helper.WithoutData{
			Code:    400,
			Message: err.Error(),
		}
		response.Response(ctx)
		return
	}

	response := &helper.WithData{
		Code:    200,
		Message: "Update",
		Data:    leave,
	}
	response.Response(ctx)
}

func (l *LeaveController) GetAllByEmployee(ctx *gin.Context) {

	var leave []models.Leave

	user, err := helper.GetUser(ctx)

	if err != nil {
		return
	}

	if err := l.DB.Find(&leave, user.Id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := &helper.WithoutData{
				Code:    400,
				Message: err.Error(),
			}
			response.Response(ctx)
			return
		}

		response := &helper.WithoutData{
			Code:    500,
			Message: err.Error(),
		}
		response.Response(ctx)
		return
	}

	response := &helper.WithData{
		Code:    200,
		Message: "success",
		Data:    leave,
	}
	response.Response(ctx)
}

func (l *LeaveController) GetOneByEmployee(ctx *gin.Context) {
	var leave models.Leave

	user, err := helper.GetUser(ctx)

	if err != nil {
		return
	}

	if err := l.DB.First(&leave, user.Id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := &helper.WithoutData{
				Code:    400,
				Message: err.Error(),
			}
			response.Response(ctx)
			return
		}
		response := &helper.WithoutData{
			Code:    500,
			Message: err.Error(),
		}
		response.Response(ctx)
		return
	}

	response := &helper.WithData{
		Code:    200,
		Message: "success",
		Data:    leave,
	}
	response.Response(ctx)
}
