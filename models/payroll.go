package models

import "time"

type Payroll struct {
	Id          int64     `gorm:"primaryKey" json:"id"`
	EmployeeId  int64     `json:"employee_id"`
	DailySalary float64   `json:"daily_salary"`
	Absence     int       `json:"absence"`
	Bonus       float64   `json:"bonus"`
	Tax         float64   `json:"tax"`
	TotalHour   float64   `json:"total_hour"`
	Total       float64   `json:"total"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
