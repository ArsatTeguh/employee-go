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

type PositionContoroller struct {
	DB *gorm.DB
}

func (p *PositionContoroller) GetAllPosition(ctx *gin.Context) {
	if valid := helper.Premission(ctx); valid != nil {
		return
	}

	var pst []models.Position

	if err := p.DB.Preload("Attedance").Find(&pst).Error; err != nil {
		helper.ErrorServer(err, ctx)
		return
	}

	res := &helper.WithData{
		Code:    200,
		Message: "Success get all positions",
		Data:    pst,
	}

	res.Response(ctx)

}

func (p *PositionContoroller) GetOnePosition(ctx *gin.Context) {
	if valid := helper.Premission(ctx); valid != nil {
		return
	}

	var pst models.Position
	id := ctx.Param("id")
	i, _ := strconv.ParseInt(id, 10, 64)
	query := p.DB.Model(&pst).Preload("Attedance")

	if err := query.Where("id = ?", i).First(&pst).Error; err != nil {
		helper.ErrorServer(err, ctx)
		return
	}

	res := helper.WithData{
		Code:    200,
		Message: "Success Get Position",
		Data:    pst,
	}

	res.Response(ctx)

}

func (p *PositionContoroller) SavePosition(ctx *gin.Context) {
	if valid := helper.Premission(ctx); valid != nil {
		return
	}

	var req []dto.RequestPosition

	if err := dto.ValidationPayload(&req, ctx); err != nil {
		return
	}

	// Save positions to the database
	// Map the request positions to models
	positions := make([]models.Position, len(req))
	for i, r := range req {
		positions[i] = r.SavePosition()
	} // Bulk create positions if err := p.DB.Create(&positions).Error; err != nil { ctx.JSON(500, gin.H{"error": err.Error()}) return }

	// Bulk create positions
	if err := p.DB.Create(&positions).Error; err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	res := &helper.WithoutData{
		Code:    201,
		Message: "Success save data",
	}
	res.Response(ctx)
}

type projectId struct {
	ProjectId int64 `json:"project_id" validate:"required"`
}

func (p *PositionContoroller) GetByProject(ctx *gin.Context) {
	var position []models.Position
	var body projectId

	if err := dto.ValidationPayload(&body, ctx); err != nil {
		return
	}

	if err := p.DB.Where("project_id = ?", body.ProjectId).Find(&position).Error; err != nil {
		helper.ErrorServer(err, ctx)
		return
	}

	res := &helper.WithData{
		Code:    200,
		Message: "Success",
		Data:    position,
	}
	res.Response(ctx)
}

func (p *PositionContoroller) Update(ctx *gin.Context) {
	if validate := helper.Premission(ctx); validate != nil {
		return
	}

	var position models.Position
	var body dto.RequestPositioUpdate
	param := ctx.Param("id")

	if param == "" {
		ctx.AbortWithStatusJSON(400, map[string]string{"message": "Required params"})
		return
	}

	id, _ := strconv.ParseInt(param, 10, 64)

	if err := dto.ValidationPayload(&body, ctx); err != nil {
		return
	}

	if err := p.DB.Model(&position).Where("id = ?", id).First(&position).Error; err != nil {
		helper.ErrorServer(err, ctx)
		return
	}

	if err := p.DB.Model(&models.Employee{}).Where("id = ?", body.EmployeeId).First(&models.Employee{}).Error; err != nil {
		helper.ErrorServer(err, ctx)
		return
	}

	if err := p.DB.Model(&models.Project{}).Where("id = ?", body.ProjectId).First(&models.Project{}).Error; err != nil {
		helper.ErrorServer(err, ctx)
		return
	}

	updateAt := time.Now()
	req := body.UpdatedPosition(&updateAt)

	if err := p.DB.Model(&position).Updates(&req).Error; err != nil {
		helper.ErrorServer(err, ctx)
		return
	}

	res := &helper.WithData{
		Code:    200,
		Message: "Update",
		Data:    position,
	}
	res.Response(ctx)
}
