package controller

import (
	"backend/dto"
	"backend/helper"
	"backend/models"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PayrollController struct {
	DB *gorm.DB
}

func (p *PayrollController) GetAll(ctx *gin.Context) {
	var payrolls models.Payroll

	if err := p.DB.Find(&payrolls).Error; err != nil {
		response := &helper.WithoutData{
			Code:    400,
			Message: "Data Kosong",
		}
		response.Response(ctx)
		return
	}
	response := &helper.WithData{
		Code:    200,
		Message: "Success Get All Payrolls",
		Data:    payrolls,
	}
	response.Response(ctx)
}

func (p *PayrollController) Payroll(ctx *gin.Context) {
	var employee models.Employee
	var attedance []models.Attedance
	var req dto.RequestPayroll

	if err := dto.ValidationPayload(&req, ctx); err != nil {
		return
	}

	if err := p.DB.Model(&employee).Preload("Wallet").Where("id = ?", req.EmployeeId).First(&employee).Error; err != nil {
		helper.ErrorServer(err, ctx)
		return
	}

	datePattern := "%" + strings.ToLower(req.Date) + "%"
	if err := p.DB.Model(&attedance).Where("lower(chekin) LIKE ? AND employee_id = ?", datePattern, req.EmployeeId).Find(&attedance).Error; err != nil {
		helper.ErrorServer(err, ctx)
		return
	}

	salary_monthly := employee.Wallet.MonthlySalary
	attedances := attedance

	result_payroll := req.CalculationPayroll(salary_monthly, attedances)

	if err := p.DB.Create(&result_payroll).Error; err != nil {
		response := &helper.WithoutData{
			Code:    400,
			Message: err.Error(),
		}
		response.Response(ctx)
		return
	}

	response := &helper.WithData{
		Code:    201,
		Message: "Payroll berhasil disimpan",
		Data:    result_payroll,
	}
	response.Response(ctx)
}

func (p *PayrollController) EmailPayslip(ctx *gin.Context) {

	file, err := ctx.FormFile("pdf")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to retrieve PDF file",
		})
		return
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to open PDF file",
		})
		return
	}
	defer src.Close()

	// Read file contents
	fileBytes, err := io.ReadAll(src)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read PDF file",
		})
		return
	}

	// Get date from form
	date := ctx.PostForm("date")

	// Use a channel for async processing if needed
	go helper.SendPayslip(fileBytes, "arsatteguh@gmail.com", date)

	response := &helper.WithoutData{
		Code:    201,
		Message: "Payroll berhasil dikirim",
	}
	response.Response(ctx)
}
