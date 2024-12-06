package helper

import (
	"backend/models"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AssosiationProject(projectId int64, ctx *gin.Context, e *gorm.DB) error {
	var project models.Project

	// Start a transaction
	err := e.Transaction(func(tx *gorm.DB) error {
		// Find the employee first
		if err := tx.Where("id = ?", projectId).First(&project).Error; err != nil {
			return err
		}

		// Delete associated records
		associatedModels := []interface{}{
			&models.Position{},
			&models.Attedance{},
			&models.Task{},
		}

		for _, model := range associatedModels {
			if err := tx.Unscoped().Where("project_id = ?", projectId).Delete(model).Error; err != nil {
				return err
			}
		}

		if err := tx.Unscoped().Delete(&project).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

func AssociationEmploee(employeeId int64, ctx *gin.Context, e *gorm.DB) error {

	var employee models.Employee

	// Start a transaction
	err := e.Transaction(func(tx *gorm.DB) error {
		// Find the employee first
		if err := tx.Where("id = ?", employeeId).First(&employee).Error; err != nil {
			return err
		}

		// Delete associated picture file if exists
		if employee.Picture != "" {
			os.Remove(employee.Picture)
		}

		// Delete associated records
		associatedModels := []interface{}{
			&models.Project{},
			&models.Wallet{},
			&models.Position{},
			&models.Payroll{},
			&models.Leave{},
			&models.Task{},
		}

		for _, model := range associatedModels {
			if err := tx.Unscoped().Where("employee_id = ?", employeeId).Delete(model).Error; err != nil {
				return err
			}
		}

		// Delete the employee
		if err := tx.Unscoped().Delete(&employee).Error; err != nil {
			return err
		}

		return nil
	})

	return err

}
