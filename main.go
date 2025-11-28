package main

import (
	"log"
	"os"

	"github.com/eetmad/backend/database"
	"github.com/eetmad/backend/models"
	"github.com/eetmad/backend/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// تحميل المتغيرات من .env
	godotenv.Load()

	// الاتصال بقاعدة البيانات
	database.Connect()

	// إنشاء الجداول تلقائيًا
	database.DB.AutoMigrate(&models.User{})

	// إعداد Gin
	r := gin.Default()

	// جميع الروتات
	api := r.Group("/api/v1")
	routes.AuthRoutes(api)

	// تشغيل السيرفر
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Eetmad Backend شغال 100%% على http://localhost:%s", port)
	log.Fatal(r.Run(":" + port))
}
