package dto

import (
	"backend/models"
	"time"
)

type RequestLeave struct {
	LeaveType string  `json:"leave_type" binding:"required"`
	StartDate string  `json:"start_date" binding:"required"`
	EndDate   string  `json:"end_date" binding:"required"`
	Status    *string `json:"status,omitempty"`
}

func (r RequestLeave) SavePosition(id int64) models.Leave {
	status := "PENDING"
	layout := "2006-01-02"
	start, err1 := time.Parse(layout, r.StartDate)
	if err1 != nil {
		panic(err1.Error())
	}

	end, err2 := time.Parse(layout, r.EndDate)
	if r.Status == nil {
		r.Status = &status
	}
	if err2 != nil {
		panic(err2.Error())
	}

	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	end = time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.UTC)

	return models.Leave{
		LeaveType:  r.LeaveType,
		StartDate:  start,
		EndDate:    end,
		Status:     r.Status,
		EmployeeId: id,
	}

}
