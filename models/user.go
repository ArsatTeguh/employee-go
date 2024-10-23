package models

type User struct {
	Id       int64    `gorm:"primaryKey" json:"id"`
	Email    string   `gorm:"type:varchar(50);not null;" json:"email"`
	Password string   `gorm:"type:varchar(255);not null;" json:"password"`
	Employee Employee `json:"-" gorm:"foreignKey:UserId"`
	Role     *string  `gorm:"type:varchar(10);" json:"role,omitempty"`
	Tokenize *string  `gorm:"type:varchar(225);" json:"tokenize,omitempty"`
}
