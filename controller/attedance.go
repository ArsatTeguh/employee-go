package controller

import (
	"backend/dto"
	"backend/helper"
	"backend/models"
	"fmt"
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

	if errs != nil {
		return
	}
	// open caching

	if err := a.DB.Where("employee_id = ?", user.Id).Find(&attedance).Error; err != nil {
		helper.ErrorServer(err, ctx)
		return
	}

	if len(attedance) == 0 {
		ctx.AbortWithStatusJSON(400, map[string]string{"message": "Data not Found"})
		return
	}

	res := helper.WithData{
		Code:    200,
		Message: "Success Get Attedances",
		Data:    attedance,
	}
	res.Response(ctx)
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

type AttedanceId struct {
	Id int64 `json:"id" binding:"required"`
}

func (a *AttedanceController) Update(ctx *gin.Context) { // attedances id body

	var chekout models.Attedance
	var att AttedanceId

	if err := dto.ValidationPayload(&att, ctx); err != nil {
		return
	}

	if _, err := a.chekinExist(ctx, &chekout, att.Id); err != nil {
		return
	}

	res := helper.WithData{
		Code:    200,
		Message: "Checkout Success",
		Data:    chekout,
	}
	res.Response(ctx)
}

func (a *AttedanceController) chekinExist(ctx *gin.Context, chekout *models.Attedance, attId int64) (models.Attedance, error) {
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

	if err := a.DB.Where("id = ? AND chekout IS NULL", attId).First(&chekout).Error; err != nil { // attedance_id
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

	if err := a.DB.Where("id = ? AND is_chekin = ?", body, 0).First(&employee).Error; err != nil {
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
