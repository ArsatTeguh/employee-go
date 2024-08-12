package controller

import (
	"backend/dto"
	"backend/helper"
	"backend/models"
	"errors"

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
	var req dto.RequestPayroll

	if body := ctx.ShouldBind(&req); body != nil {
		response := &helper.WithoutData{
			Code:    400,
			Message: body.Error(),
		}
		response.Response(ctx)
		return
	}

	query := p.DB.Model(&employee).Preload("Wallet").Preload("Position").Preload("Position.Attedance")
	if err := query.Where("id = ?", req.EmployeeId).First(&employee).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := &helper.WithoutData{
				Code:    400,
				Message: err.Error(),
			}
			response.Response(ctx)
			return
		}
	}

	if len(employee.Position) == 0 {
		response := &helper.WithoutData{
			Code:    400,
			Message: "Data Position tidak ada",
		}
		response.Response(ctx)
		return
	}
	salary_monthly := employee.Wallet.MonthlySalary
	attedances := employee.Position[0].Attedance

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
