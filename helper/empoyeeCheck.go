package helper

import (
	"backend/models"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func EmployeeExist(id int64, db *gorm.DB) (models.Employee, error) {
	var employee models.Employee
	query := db.Model(&employee).Preload("Wallet").Preload("Position").Preload("Position.Attedance").Preload("Leave")
	if err := query.Where("id = ?", id).First(&employee).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return employee, fmt.Errorf(err.Error())
		}
		return employee, fmt.Errorf(err.Error())
	}
	return employee, nil
}
