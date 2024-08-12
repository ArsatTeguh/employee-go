package dto

import (
	"backend/helper"
	"backend/models"
)

type RequestPayroll struct {
	Absen      *int     `json:"absen"`
	EmployeeId int64    `json:"employee_id" binding:"required"`
	Bonus      *float64 `json:"bonus"`
	Tax        *float64 `json:"tax"`
	Status     string   `json:"status" binding:"required"`
}

func (c RequestPayroll) CalculationPayroll(salary_monthly float64, position []models.Attedance) models.Payroll {
	zero := 0
	var tax_zero float64 = 0
	if c.Tax == nil {
		c.Tax = &tax_zero
	}
	if c.Bonus == nil {
		c.Bonus = &tax_zero
	}
	if c.Absen == nil {
		c.Absen = &zero
	}

	salary_day := salary_monthly / 30
	subtraction := salary_day * float64(*c.Absen)

	tax_decimal := *c.Tax / 100                                    // convert pajak ke decimal
	tax := (salary_monthly + *c.Bonus - subtraction) * tax_decimal // hitung pajak berdasarkan jumlah salary
	total := (salary_monthly + *c.Bonus - subtraction) - tax
	total_hourse := helper.CalculationWorkMonthly(position, c.EmployeeId)

	return models.Payroll{
		EmployeeId:  c.EmployeeId,
		DailySalary: salary_day,
		Absence:     *c.Absen,
		Bonus:       *c.Bonus,
		Tax:         *c.Tax,
		Total:       total,
		Status:      c.Status,
		TotalHour:   total_hourse,
	}
}
