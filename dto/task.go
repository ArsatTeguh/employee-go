package dto

import (
	"backend/models"
	"time"
)

type CreateTask struct {
	ProjectId   int64   `json:"project_id" validate:"required"`
	Title       string  `json:"title" validate:"required"`
	ActionName  *string `json:"action_name,omitempty"`
	Description string  `json:"description" validate:"required"`
	Status      int64   `json:"status" validate:"required"`
	Level       int64   `json:"level" validate:"required"`
	Action      int64   `json:"action" validate:"required"`
	EmployeeId  int64   `json:"employee_id" validate:"required"`
}

type UpdateTask struct {
	ActionName  *string    `json:"action_name"`
	Description string     `json:"description"`
	Title       string     `json:"title"`
	Status      int64      `json:"status"`
	Level       int64      `json:"level"`
	Action      int64      `json:"action"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

func (t *CreateTask) PayloadTask() models.Task {

	return models.Task{
		ProjectId:   t.ProjectId,
		Description: t.Description,
		Status:      t.Status,
		Level:       t.Level,
		Action:      t.Action,
		ActionName:  t.ActionName,
		EmployeeId:  t.EmployeeId,
		Title:       t.Title,
	}

}
func (t *UpdateTask) PayloadTaskUpdate(updatedTime *time.Time) models.Task {

	save := models.Task{
		Description: t.Description,
		Status:      t.Status,
		Level:       t.Level,
		Action:      t.Action,
		ActionName:  t.ActionName,
		Title:       t.Title,
	}
	if updatedTime != nil {
		save.UpdatedAt = *updatedTime
	}

	return save
}
