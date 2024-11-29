package dto

type UpdateWallet struct {
	HourlySalary  int64  `json:"hourly_salary" validate:"required"`
	MonthlySalary int64  `json:"monthly_salary" validate:"required"`
	NoRekening    int64  `json:"no_rekening" validate:"required"`
	NameBanking   string `json:"name_banking" validate:"required"`
	TypeBanking   string `json:"type_banking" validate:"required"`
}
