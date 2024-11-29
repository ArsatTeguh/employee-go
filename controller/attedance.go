package controller

import (
	"backend/dto"
	"backend/helper"
	"backend/models"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AttedanceController struct {
	DB *gorm.DB
}

func (a *AttedanceController) GetAll(ctx *gin.Context) {
	var attedance []models.Attedance
	user, errs := helper.GetUser(ctx)

	search := ctx.Query("date")

	var totalCount int64
	if errs != nil {
		return
	}
	// open caching

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	sizePage, _ := strconv.Atoi((ctx.Copy().DefaultQuery("sizePage", "10")))
	offset := (page - 1) * sizePage

	query := a.DB.Model(&attedance).Preload("Project").Where("employee_id = ?", user.Id)

	if matched, _ := regexp.MatchString(`^\d{4}-\d{2}$`, search); !matched {
		response := &helper.WithoutData{
			Code:    400,
			Message: "invalid date format. Use YYYY-MM",
		}
		response.Response(ctx)
		return
	}

	if search != "" {
		query = query.Where("DATE_FORMAT(chekin, '%Y-%m') = ?", search)
	}

	query.Count(&totalCount)
	query.Offset(offset).Limit(sizePage).Find(&attedance)

	if len(attedance) == 0 {
		response := &helper.WithoutData{
			Code:    400,
			Message: "Data empty",
		}
		response.Response(ctx)
		return
	}

	response := &helper.WithData{
		Code:    200,
		Message: "Success Get Attedances",
		Data: map[string]any{
			"attedances": attedance,                                            // data
			"totalAll":   totalCount,                                           // total data all page
			"total":      len(attedance),                                       // total data per page
			"page":       page,                                                 // current page
			"sizePage":   sizePage,                                             // maximum data per page
			"totalPages": (totalCount + int64(sizePage) - 1) / int64(sizePage), // total all page
		},
	}
	response.Response(ctx)
}

func (a *AttedanceController) GetOne(ctx *gin.Context) {
	var attedance models.Attedance
	param := ctx.Param("id") // employee id
	id, _ := strconv.Atoi(param)
	validate := helper.Premission(ctx)

	if validate != nil {
		return
	}

	if err := a.DB.Where("employee_id = ?", id).Find(&attedance).Error; err != nil {
		helper.ErrorServer(err, ctx)
		return
	}

	res := helper.WithData{
		Code:    200,
		Message: "Success get one Attedances",
		Data:    attedance,
	}
	res.Response(ctx)

}

func (a *AttedanceController) Created(ctx *gin.Context) {
	user, errs := helper.GetUser(ctx)

	if errs != nil {
		return
	}

	var attedance dto.RequestAttedance

	if err := dto.ValidationPayload(&attedance, ctx); err != nil {
		return
	}

	zero := a.zeroChekin(user.Id, 0)

	if zero != nil {
		res := helper.WithoutData{
			Code:    400,
			Message: "User dalam keadaan chekin",
		}
		res.Response(ctx)
		return
	}

	currentTime := time.Now()

	response := attedance.SavePosition(&currentTime, nil, user.Id)

	if isInsert := a.DB.Create(&response).Error; isInsert != nil {
		res := helper.WithoutData{
			Code:    500,
			Message: isInsert.Error(),
		}
		res.Response(ctx)
		return
	}

	if updateEmployee := a.DB.Model(&models.Employee{}).Where("id = ?", user.Id).Update("is_chekin", 1).Error; updateEmployee != nil {
		res := helper.WithoutData{
			Code:    500,
			Message: "is_chekin employee tidak berhasil di update",
		}
		res.Response(ctx)
		return
	}

	res := helper.WithData{
		Code:    201,
		Message: "Checkin Success",
		Data:    response,
	}
	res.Response(ctx)

}

func (a *AttedanceController) Update(ctx *gin.Context) { // attedances id body

	var chekout models.Attedance

	if _, err := a.chekinExist(ctx, &chekout); err != nil {
		return
	}

	res := helper.WithData{
		Code:    200,
		Message: "Checkout Success",
		Data:    chekout,
	}
	res.Response(ctx)
}

func (a *AttedanceController) chekinExist(ctx *gin.Context, chekout *models.Attedance) (models.Attedance, error) {
	var employee models.Employee
	user, err := helper.GetUser(ctx)

	if err != nil {
		res := helper.WithoutData{
			Code:    203,
			Message: "No match users",
		}
		res.Response(ctx)
		return *chekout, fmt.Errorf(err.Error())
	}

	if err := a.DB.Where("employee_id = ? AND chekout IS NULL", user.Id).First(&chekout).Error; err != nil { // userId
		res := helper.WithoutData{
			Code:    400,
			Message: "Attedance tidak ditemukan atau user sudah checkout",
		}
		res.Response(ctx)
		return *chekout, fmt.Errorf(err.Error())
	}

	if user.Id != chekout.EmployeeId {
		res := helper.WithoutData{
			Code:    400,
			Message: "checkout dapat dilakukan dengan akun yang login saja",
		}
		res.Response(ctx)
		return *chekout, fmt.Errorf(err.Error())
	}

	zero := a.zeroChekin(chekout.EmployeeId, 1)

	if zero != nil {
		return *chekout, fmt.Errorf("Error")
	}

	if err := a.DB.First(&employee, "id = ?", chekout.EmployeeId).Error; err != nil {
		helper.ErrorServer(err, ctx)
		return *chekout, fmt.Errorf(err.Error())
	}

	currentTime := time.Now()
	chekout.Chekout = &currentTime
	working := helper.CalculationWorkHour(*chekout.Chekin, currentTime)
	chekout.Working_house = working
	chekout.UpdatedAt = currentTime

	if isUpdate := a.DB.Updates(&chekout).Error; isUpdate != nil {
		helper.ErrorServer(isUpdate, ctx)
		return *chekout, fmt.Errorf(isUpdate.Error())

	}

	condition := false

	employee.IsChekin = &condition
	if err := a.DB.Updates(&employee).Error; err != nil {
		res := helper.WithoutData{
			Code:    500,
			Message: "Employee isCheckout tidak ter update",
		}
		res.Response(ctx)
		return *chekout, fmt.Errorf(err.Error())
	}

	return *chekout, nil
}

type BodyId struct {
	Employee_id int64 `json:"employee_id" validate:"required"`
}

func (a *AttedanceController) EmployeeCheckout(ctx *gin.Context) { // body id employe
	var employee models.Employee
	var body BodyId

	if err := dto.ValidationPayload(&body, ctx); err != nil {
		return
	}

	if err := a.DB.Where("id = ? AND is_chekin = ?", body.Employee_id, 0).First(&employee).Error; err != nil {
		helper.ErrorServer(err, ctx)
		return
	}

	res := helper.WithoutData{
		Code:    200,
		Message: "TRUE",
	}
	res.Response(ctx)

}

func (a *AttedanceController) zeroChekin(id int64, num int) error {
	if err := a.DB.Where("id = ? AND is_chekin = ?", id, num).First(&models.Employee{}).Error; err != nil {
		return fmt.Errorf(err.Error())
	}
	return nil
}
