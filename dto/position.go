package dto

import (
	"backend/models"
	"time"
)

type RequestPosition struct {
	Position   string `json:"position" binding:"required"`
	Status     string `json:"status" binding:"required"`
	EmployeeId int64  `json:"employee_id" binding:"required"`
	ProjectId  int64  `json:"project_id" binding:"required"`
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
	Position   string     `json:"position" binding:"required"`
	Status     string     `json:"status" binding:"required"`
	EmployeeId int64      `json:"employee_id" binding:"required"`
	ProjectId  int64      `json:"project_id" binding:"required"`
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
