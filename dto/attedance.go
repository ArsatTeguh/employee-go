package dto

import (
	"backend/models"
	"time"
)

type RequestAttedance struct {
	Location   string `json:"location" validate:"required"`
	EmployeeId *int64 `json:"employee_id,omitempty"`
	ProjectId  int64  `json:"project_id" validate:"required"`
	PositionId int64  `json:"position_id" validate:"required"`
}

func (c RequestAttedance) SavePosition(chekin *time.Time, chekout *time.Time, id int64) models.Attedance {

	attedance := models.Attedance{
		Location:   c.Location,
		EmployeeId: id,
		ProjectId:  c.ProjectId,
		PositionId: c.PositionId,
	}

	attedance.Chekin = chekin

	return attedance
}
