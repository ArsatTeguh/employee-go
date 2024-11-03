package dto

import (
	"backend/models"
	"time"
)

type RequestPosition struct {
	Position   string `json:"position" validate:"required"`
	Status     string `json:"status" validate:"required"`
	EmployeeId int64  `json:"employee_id" validate:"required"`
	ProjectId  int64  `json:"project_id" validate:"required"`
}

func (c RequestPosition) SavePosition() models.Position {

	return models.Position{
		Position:   c.Position,
		Status:     c.Status,
		EmployeeId: c.EmployeeId,
		ProjectId:  c.ProjectId,
	}
}

type RequestPositioUpdate struct {
	Position   string     `json:"position" validate:"required"`
	Status     string     `json:"status" validate:"required"`
	EmployeeId int64      `json:"employee_id" validate:"required"`
	ProjectId  int64      `json:"project_id" validate:"required"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
}

func (c RequestPositioUpdate) UpdatedPosition(updateTime *time.Time) models.Position {

	return models.Position{
		Position:   c.Position,
		Status:     c.Status,
		EmployeeId: c.EmployeeId,
		ProjectId:  c.ProjectId,
		UpdatedAt:  *updateTime,
	}
}
