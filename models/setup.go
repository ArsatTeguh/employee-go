package models

import (
	"database/sql"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func ConnectionDb() {
	sqlDB, _ := sql.Open("mysql", GetEnv("MYSQL", "env MYSQL kosong"))
	database, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	// database, err := gorm.Open(mysql.Open(GetEnv("MYSQL", "env MYSQL kosong")))

	if err != nil {
		panic(err)
	}
	database.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&User{},
		&Project{},
		&Employee{},
		&Wallet{},
		&Position{},
		&Attedance{},
		&Payroll{},
		&Leave{},
	)
	DB = database

}
