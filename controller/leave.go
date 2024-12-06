package controller

import (
	"backend/dto"
	"backend/helper"
	"backend/models"
	"errors"
	"strconv"
	"strings"
	"time"

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
	user, _ := helper.GetUser(ctx)

	if err := dto.ValidationPayload(&body, ctx); err != nil {
		return
	}

	res := body.SavePosition(user.Id)

	if err := l.DB.Create(&res).Error; err != nil {
		response := &helper.WithoutData{
			Code:    400,
			Message: err.Error(),
		}
		response.Response(ctx)
		return
	}

	data := helper.FormatEmailLeave{
		LeaveType: res.LeaveType,
		StartDate: body.StartDate,
		EndDate:   body.EndDate,
		Status:    res.Status,
		Employee:  user.Email,
	}

	go helper.SendEmail(data)

	response := &helper.WithData{
		Code:    201,
		Message: "Insert",
		Data:    body,
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

	if err := l.DB.Where("id = ? AND employee_id = ?", id, body.EmployeeId).First(&leave).Error; err != nil {
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

type PayloadByIdEmployee struct {
	EmployeeId int64 `json:"employee_id" binding:"required"`
}

type ResponseLeave struct {
	Id         int64     `json:"id" `
	EmployeeId int64     `json:"employee_id"`
	LeaveType  string    `json:"leave_type"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

func (l *LeaveController) GetAllByEmployee(ctx *gin.Context) {

	var leave []models.Leave
	var result []ResponseLeave
	param := ctx.Param("id")

	if param == "" {
		ctx.AbortWithStatusJSON(400, map[string]string{"message": "Param Required"})
		return
	}

	id, _ := strconv.ParseInt(param, 10, 64)

	if err := l.DB.Model(&leave).Where("employee_id = ?", id).Find(&result).Error; err != nil {
		helper.ErrorServer(err, ctx)
		return
	}

	if len(result) == 0 {
		response := &helper.WithoutData{
			Code:    400,
			Message: "Data Not Found",
		}
		response.Response(ctx)
		return
	}

	response := &helper.WithData{
		Code:    200,
		Message: "success",
		Data:    result,
	}
	response.Response(ctx)
}

func (l *LeaveController) GetOneByEmployee(ctx *gin.Context) {
	var leave []models.Leave
	var res []ResponseLeave
	user, err := helper.GetUser(ctx)

	if err != nil {
		return
	}

	if err := l.DB.Model(&leave).Where("employee_id = ? ", user.Id).Find(&res).Error; err != nil {
		helper.ErrorServer(err, ctx)
		return
	}

	if len(res) == 0 {
		response := &helper.WithoutData{
			Code:    400,
			Message: "Data Not Found",
		}
		response.Response(ctx)
		return
	}
	response := &helper.WithData{
		Code:    200,
		Message: "success",
		Data:    res,
	}
	response.Response(ctx)
}

func (l *LeaveController) Delete(ctx *gin.Context) {
	err := helper.Premission(ctx)
	if err != nil {
		return
	}

	var leave models.Leave
	param := ctx.Param("id") // leave id
	id, _ := strconv.Atoi(param)

	if err := l.DB.Delete(&leave, id).Error; err != nil {
		helper.ErrorServer(err, ctx)
		return
	}

	response := &helper.WithoutData{
		Code:    200,
		Message: "success",
	}
	response.Response(ctx)
}
