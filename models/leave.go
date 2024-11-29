package models

import (
	"encoding/json"
	"time"
)

type Leave struct {
	Id         int64     `json:"id" gorm:"primaryKey"`
	EmployeeId int64     `gorm:"index" json:"employee_id"`
	LeaveType  string    `gorm:"size:50" json:"leave_type"`
	StartDate  time.Time `gorm:"size:50" json:"start_date"`
	EndDate    time.Time `gorm:"size:50" json:"end_date"`
	Employee   Employee  `json:"employee"`
	Status     *string   `gorm:"size:50" json:"status,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type LeaveFormat struct {
	employee   string
	leave_type string
	start_date time.Time
	end_date   time.Time
	status     string
}

func (l Leave) MarshalJSON() ([]byte, error) {
	type Alias Leave
	return json.Marshal(&struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		Alias
	}{
		StartDate: l.StartDate.Format("2006-01-02"),
		EndDate:   l.EndDate.Format("2006-01-02"),
		Alias:     (Alias)(l),
	})
}
