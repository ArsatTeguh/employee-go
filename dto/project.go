package dto

import (
	"backend/models"
	"time"
)

type RequestProject struct {
	Name       string     `json:"name" validate:"required"`
	Estimation string     `json:"estimation" validate:"required"`
	Status     string     `json:"status" validate:"required"`
	EmployeeId *int64     `json:"employee_id"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

func (c RequestProject) SaveProject(updatedTime *time.Time, employeId int64) models.Project {

	res := models.Project{
		Name:       c.Name,
		Estimation: c.Estimation,
		Status:     c.Status,
		EmployeeId: employeId,
	}

	if updatedTime != nil {
		res.UpdatedAt = *updatedTime
	}

	return res
}
