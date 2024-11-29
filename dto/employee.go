package dto

import (
	"backend/models"
)

type CreatePayload struct {
	Phone         int     `json:"phone" validate:"required"`
	Picture       string  `json:"picture,omitempty"`
	HourlySalary  float64 `json:"hourly_salary" validate:"required"`
	MonthlySalary float64 `json:"monthly_salary" validate:"required"`
	NoRekening    int     `json:"no_rekening" validate:"required"`
	NameBanking   string  `json:"name_banking" validate:"required"`
	TypeBanking   string  `json:"type_banking" validate:"required"`
}

type UpdateEmployeeStruct struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   int    `json:"phone"`
}

func (c *CreatePayload) PayloadEmployee() models.Employee {

	return models.Employee{
		Phone:   c.Phone,
		Picture: c.Picture,
	}
}

func (c *UpdateEmployeeStruct) UpdateEmployee() models.Employee {

	return models.Employee{
		Phone:   c.Phone,
		Address: c.Address,
		Name:    c.Name,
	}
}

func (c *CreatePayload) PayloadWallet(IdEmployee int64) models.Wallet {

	return models.Wallet{
		HourlySalary:  c.HourlySalary,
		MonthlySalary: c.MonthlySalary,
		NoRekening:    c.NoRekening,
		NameBanking:   c.NameBanking,
		TypeBanking:   c.TypeBanking,
		EmployeeId:    IdEmployee,
	}
}
