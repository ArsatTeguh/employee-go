package models

import "time"

type Position struct {
	Id         int64     `json:"id" gorm:"primaryKey"`
	Position   string    `json:"position" gorm:"size:50;" `
	Status     string    `json:"status" gorm:"size:50"`
	EmployeeId int64     `gorm:"index" json:"employee_id"`
	ProjectId  int64     `gorm:"index" json:"project_id"`
	Project    Project   `json:"project"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
