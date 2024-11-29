package models

import "time"

type Attedance struct {
	Id            int64      `json:"id" gorm:"primaryKey"`
	Chekin        *time.Time `json:"chekin,omitempty" gorm:"size:50;"`
	Chekout       *time.Time `json:"chekout,omitempty" gorm:"size:50;"`
	Location      string     `json:"location" gorm:"size:50"`
	Working_house float64    `json:"working_house" gorm:"size:50"`
	EmployeeId    int64      `gorm:"index" json:"employee_id"`
	ProjectId     int64      `gorm:"index" json:"project_id"`
	Project       Project    `json:"project"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
