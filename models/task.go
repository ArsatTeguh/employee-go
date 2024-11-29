package models

import "time"

type Task struct {
	Id          int64     `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title"`
	ProjectId   int64     `json:"project_id"`
	EmployeeId  int64     `json:"employee_id"`
	Description string    `json:"description"`
	Status      int64     `json:"status" gorm:"size:5;"`
	Level       int64     `json:"level" gorm:"size:5;"`
	Action      int64     `json:"action" gorm:"size:5;"`
	ActionName  *string   `json:"action_name,omitempty" gorm:"size:50;"`
	Employee    Employee  `json:"employee"`
	Project     Project   `json:"project"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
