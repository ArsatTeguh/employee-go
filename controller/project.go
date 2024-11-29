package controller

import (
	"backend/dto"
	"backend/helper"
	"backend/models"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProjectController struct {
	DB *gorm.DB
}

func (p *ProjectController) GetAllProject(ctx *gin.Context) {
	// valid := helper.Premission(ctx)

	// if valid != nil {
	// 	return
	// }

	var project []models.Project
	search := ctx.Query("name")
	var totalCount int64

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	sizePage, _ := strconv.Atoi((ctx.Copy().DefaultQuery("sizePage", "5")))
	offset := (page - 1) * sizePage

	query := p.DB.Model(&project).Preload("Position")

	// if query have some string and then
	if search != "" {
		namePattern := "%" + strings.ToLower(search) + "%"
		query = query.Where("lower(name) LIKE ?", namePattern)
	}

	query.Count(&totalCount)
	query.Offset(offset).Limit(sizePage).Find(&project)

	if len(project) == 0 {
		response := &helper.WithoutData{
			Code:    400,
			Message: "Data empty",
		}
		response.Response(ctx)
		return
	}

	response := &helper.WithData{
		Code:    200,
		Message: "Success Get Projects",
		Data: map[string]any{
			"projects":   project,                                              // data
			"totalAll":   totalCount,                                           // total data all page
			"total":      len(project),                                         // total data per page
			"page":       page,                                                 // current page
			"sizePage":   sizePage,                                             // maximum data per page
			"totalPages": (totalCount + int64(sizePage) - 1) / int64(sizePage), // total all page
		},
	}
	response.Response(ctx)
}

func (p *ProjectController) GetOne(ctx *gin.Context) {
	var project models.Project

	id := ctx.Param("id")

	query := p.DB.Model(&project).Preload("Attedance").Preload("Position")
	if err := query.Where("id = ?", id).First(&project).Error; err != nil {
		helper.ErrorServer(err, ctx)
		return
	}

	res := &helper.WithData{
		Code:    200,
		Message: "Success Get Project",
		Data:    project,
	}
	res.Response(ctx)
}

func (p *ProjectController) Saved(ctx *gin.Context) {
	if valid := helper.Premission(ctx); valid != nil {
		return
	}
	user, _ := helper.GetUser(ctx)

	var req dto.RequestProject

	if err := dto.ValidationPayload(&req, ctx); err != nil {
		return
	}

	project := req.SaveProject(nil, user.Id)

	if err := p.DB.Create(&project).Error; err != nil {
		res := &helper.WithoutData{
			Code:    500,
			Message: err.Error(),
		}
		res.Response(ctx)
		return
	}

	res := &helper.WithData{
		Code:    201,
		Message: "insert Project",
		Data:    project,
	}
	res.Response(ctx)
}

func (p *ProjectController) Update(ctx *gin.Context) {
	if valid := helper.Premission(ctx); valid != nil {
		return
	}
	user, _ := helper.GetUser(ctx)
	var project models.Project
	var body dto.RequestProject
	param := ctx.Param("id")

	if param == "" {
		ctx.AbortWithStatusJSON(400, map[string]string{"message": "Required Params"})
		return
	}

	id, _ := strconv.ParseInt(param, 10, 64)

	if err := dto.ValidationPayload(&body, ctx); err != nil {
		return
	}
	timeUpdate := time.Now()

	req := body.SaveProject(&timeUpdate, user.Id)

	qr := p.DB.Model(&project).Where("id = ?", id).First(&project)
	if err := qr.Updates(&req).Error; err != nil {
		helper.ErrorServer(err, ctx)
		return
	}

	res := &helper.WithData{
		Code:    200,
		Message: "update Project",
		Data:    project,
	}
	res.Response(ctx)
}

func (p *ProjectController) Delete(ctx *gin.Context) {
	if valid := helper.Premission(ctx); valid != nil {
		return
	}

	var project models.Project
	param := ctx.Param("id")

	if param == "" {
		ctx.AbortWithStatusJSON(400, map[string]string{"message": "Param Required"})
		return
	}

	id, _ := strconv.ParseInt(param, 10, 64)

	if err := p.DB.Where("id = ?", id).First(&project).Error; err != nil {
		helper.ErrorServer(err, ctx)
		return
	}

	// Update foreign key to nullify the dependency
	if err := p.DB.Where("project_id = ?", id).Delete(&models.Position{}).Error; err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if err := p.DB.Where("project_id = ?", id).Delete(&models.Attedance{}).Error; err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if err := p.DB.Where("project_id = ?", id).Delete(&models.Task{}).Error; err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if err := p.DB.Delete(&project, id).Error; err != nil {
		helper.ErrorServer(err, ctx)
		return
	}

	res := &helper.WithoutData{
		Code:    200,
		Message: "Delete Success",
	}
	res.Response(ctx)
}

type ProjectMasterResponse struct {
	Projects     []models.Project  `json:"projects"`
	StatusCounts []StatusCount     `json:"status_counts"`
	Employee     []TypeNewEmployee `json:"employee"`
}

type StatusCount struct {
	Status     string `json:"status"`
	TotalCount int    `json:"total_count"`
}

type TypeNewEmployee struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Picture string `json:"picture"`
	Email   string `json:"email"`
}

func (p *ProjectController) ProjectMaster(ctx *gin.Context) {
	valid := helper.Premission(ctx)
	if valid != nil {
		return
	}

	var projects []models.Project
	var statusCounts []StatusCount
	var employees []TypeNewEmployee

	// Fetch latest projects
	if err := p.DB.Model(&projects).
		Preload("Position").
		Order("created_at desc").
		Limit(3).
		Find(&projects).Error; err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		// jika terjadi error lanjutkan saja
	}

	// Fetch status counts
	if err := p.DB.Table("projects"). // Adjust table name
						Select("status, COUNT(*) as total_count").
						Group("status").
						Find(&statusCounts).Error; err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		// jika terjadi error lanjutkan saja
	}

	if err := p.DB.Table("employees").Select("name", "status", "picture", "email").Find(&employees).Error; err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
	}

	// Combine results
	res := &helper.WithData{
		Code:    200,
		Message: "Success Get Projects",
		Data: ProjectMasterResponse{
			Projects:     projects,
			StatusCounts: statusCounts,
			Employee:     employees,
		},
	}

	res.Response(ctx)
}
