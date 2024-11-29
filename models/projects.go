package models

import "time"

type Project struct {
	Id         int64       `json:"id" gorm:"primaryKey"`
	Name       string      `json:"name" gorm:"size:50;"`
	Estimation string      `json:"Estimation" gorm:"size:50;"`
	Status     string      `json:"status" gorm:"size:50;"`
	Position   []Position  `json:"position" gorm:"foreignKey:ProjectId"`
	Attedance  []Attedance `json:"attedance" gorm:"foreignKey:ProjectId"`
	Task       []Task      `json:"task" gorm:"foreignKey:ProjectId"`
	EmployeeId int64       `gorm:"index" json:"employee_id"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}
