package controller

import (
	"backend/dto"
	"backend/helper"
	"backend/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskController struct {
	DB *gorm.DB
}

func (t *TaskController) GetOne(ctx *gin.Context) {
	var task []models.Task

	param := ctx.Param("projectId")

	if param == "" {
		ctx.AbortWithStatusJSON(400, map[string]string{"message": "Required Params"})
		return
	}

	id, _ := strconv.ParseInt(param, 10, 64)

	if err := t.DB.Preload("Employee").Preload("Project").
		Where("project_id = ?", id).Find(&task).Error; err != nil {
		return
	}

	if len(task) == 0 {
		res := &helper.WithoutData{
			Code:    404,
			Message: "Not Record",
		}
		res.Response(ctx)
		return
	}

	taskDetails := []map[string]interface{}{}
	for _, task := range task {
		taskDetails = append(taskDetails, map[string]interface{}{
			"id":          task.Id,
			"title":       task.Title,
			"action_name": task.ActionName,
			"description": task.Description,
			"status":      task.Status,
			"level":       task.Level,
			"action":      task.Action,
			"employee":    task.Employee.Name,
			"project":     task.Project.Name,
			"created_at":  task.CreatedAt,
		})
	}

	res := &helper.WithData{
		Code:    200,
		Message: "Success Get Task",
		Data:    taskDetails,
	}
	res.Response(ctx)

}

func (t *TaskController) SaveTask(ctx *gin.Context) {

	var req dto.CreateTask

	if err := dto.ValidationPayload(&req, ctx); err != nil {
		return
	}

	body := req.PayloadTask()

	if err := t.DB.Create(&body).Error; err != nil {
		res := &helper.WithoutData{
			Code:    500,
			Message: err.Error(),
		}
		res.Response(ctx)
		return
	}

	if err := t.DB.Preload("Employee").Preload("Project").
		Where("id", body.Id).First(&body).Error; err != nil {
		return
	}

	taskDetails := map[string]interface{}{
		"id":          body.Id,
		"title":       body.Title,
		"action_name": body.ActionName,
		"description": body.Description,
		"status":      body.Status,
		"level":       body.Level,
		"action":      body.Action,
		"employee":    body.Employee.Name,
		"project":     body.Project.Name,
		"created_at":  body.CreatedAt,
	}

	res := &helper.WithData{
		Code:    201,
		Message: "insert Task",
		Data:    taskDetails,
	}
	res.Response(ctx)

}

func (t *TaskController) Update(ctx *gin.Context) {

	param := ctx.Param("id")

	if param == "" {
		ctx.AbortWithStatusJSON(400, map[string]string{"message": "Required Params"})
		return
	}

	var req dto.UpdateTask
	var task models.Task
	id, _ := strconv.ParseInt(param, 10, 64)

	if err := dto.ValidationPayload(&req, ctx); err != nil {
		return
	}
	timeUpdate := time.Now()

	body := req.PayloadTaskUpdate(&timeUpdate)

	qr := t.DB.Model(&task).Where("id = ?", id).First(&task)
	if err := qr.Updates(&body).Error; err != nil {
		helper.ErrorServer(err, ctx)
		return
	}

	res := &helper.WithData{
		Code:    200,
		Message: "update Task",
		Data:    task,
	}
	res.Response(ctx)

}

func (t *TaskController) Delete(ctx *gin.Context) {
	param := ctx.Param("id")

	if param == "" {
		ctx.AbortWithStatusJSON(400, map[string]string{"message": "Required Params"})
		return
	}
	id, _ := strconv.ParseInt(param, 10, 64)

	var task models.Task

	if err := t.DB.Where("id = ?", id).First(&task).Error; err != nil {
		helper.ErrorServer(err, ctx)
		return
	}

	if err := t.DB.Delete(&task, id).Error; err != nil {
		helper.ErrorServer(err, ctx)
		return
	}

	res := &helper.WithoutData{
		Code:    200,
		Message: "Success Delete",
	}
	res.Response(ctx)
}
