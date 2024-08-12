package models

import "time"

type Wallet struct {
	Id            int64     `json:"id" gorm:"primaryKey"`
	HourlySalary  float64   `json:"hourly_salary" gorm:"size:50;"`
	MonthlySalary float64   `json:"monthly_salary" gorm:"size:50;"`
	NoRekening    int       `json:"no_rekening" gorm:"size:50;"`
	NameBanking   string    `json:"name_banking" gorm:"size:50;"`
	TypeBanking   string    `json:"type_banking" gorm:"size:50;"`
	EmployeeId    int64     `gorm:"index" json:"employee_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
