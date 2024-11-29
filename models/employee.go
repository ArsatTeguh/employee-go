package models

import "time"

type Employee struct {
	Id        int64      `gorm:"primaryKey" json:"id"`
	Name      string     `gorm:"size:50; not null;" json:"name"`
	Address   string     `gorm:"size:100; not null;" json:"address"`
	Status    string     `gorm:"size:20; not null;" json:"status"`
	Phone     int        `gorm:"size:50;" json:"phone"`
	Email     string     `gorm:"size:100;" json:"email"`
	Picture   string     `json:"Picture"`
	Project   []Project  `json:"project" gorm:"foreignKey:EmployeeId"`
	Position  []Position `json:"position" gorm:"foreignKey:EmployeeId"`
	Wallet    Wallet     `json:"wallet" gorm:"foreignKey:EmployeeId"`
	Leave     []Leave    `json:"leave" gorm:"foreignKey:EmployeeId"`
	Tasks     []Task     `json:"task" gorm:"foreignKey:EmployeeId"`
	IsChekin  *bool      `json:"isChekin" gorm:"type:boolean;default:false"`
	UserId    int64      `gorm:"index" json:"user_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
