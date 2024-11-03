package main

import (
	"backend/helper"
	"backend/models"
	"backend/routers"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error Running to .env file")
	}

	models.ConnectionDb()
}

func main() {
	server := gin.Default()
	server.Static("/assets", "./assets")
	server.MaxMultipartMemory = 1024 * 1024 * 1 // 1MB
	databases := models.DB
	server.Use(helper.CORSMiddleware())

	r := routers.SetupRouter(server, databases)
	r.Run()
}
