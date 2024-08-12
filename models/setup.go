package models

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func ConnectionDb() {
	database, err := gorm.Open(mysql.Open(GetEnv("MYSQL", "env MYSQL kosong")), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

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
