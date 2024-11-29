package controller

import (
	"backend/dto"
	"backend/helper"
	"backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WalletContoroller struct {
	DB *gorm.DB
}

func (p *WalletContoroller) GetOneWallet(ctx *gin.Context) {
	if valid := helper.Premission(ctx); valid != nil {
		return
	}

	var wall models.Wallet

	id := ctx.Param("id")
	i, _ := strconv.ParseInt(id, 10, 64)
	query := p.DB.Model(&wall)

	if err := query.Where("employee_id = ?", i).First(&wall).Error; err != nil {
		helper.ErrorServer(err, ctx)
		return
	}

	res := helper.WithData{
		Code:    200,
		Message: "Success Get Wallet",
		Data:    wall,
	}

	res.Response(ctx)

}

func (w *WalletContoroller) UpdatedWallet(ctx *gin.Context) {
	if valid := helper.Premission(ctx); valid != nil {
		return
	}

	id := ctx.Param("id")
	i, _ := strconv.ParseInt(id, 10, 64)

	var payload dto.UpdateWallet
	var wallet models.Wallet

	if err := dto.ValidationPayload(&payload, ctx); err != nil {
		return
	}

	if err := w.DB.Model(&wallet).Where("employee_id = ?", i).Updates(&payload).Error; err != nil {
		helper.ErrorServer(err, ctx)
		return
	}

	res := helper.WithData{
		Code:    200,
		Message: "Update Wallet",
		Data:    payload,
	}

	res.Response(ctx)
}
