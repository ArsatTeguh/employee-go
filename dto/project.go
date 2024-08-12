package dto

import (
	"backend/models"
	"time"
)

type RequestProject struct {
	Name       string     `json:"name" binding:"required"`
	Estimation string     `json:"estimation" binding:"required"`
	Status     string     `json:"status" binding:"required"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

func (c RequestProject) SaveProject(updatedTime *time.Time) models.Project {

	res := models.Project{
		Name:       c.Name,
		Estimation: c.Estimation,
		Status:     c.Status,
	}

	if updatedTime != nil {
		res.UpdatedAt = *updatedTime
	}

	return res
}
