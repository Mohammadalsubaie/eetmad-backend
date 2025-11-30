package main

import (
	"log"
	"os"

	"github.com/eetmad/backend/database"
	"github.com/eetmad/backend/routes"
	"github.com/eetmad/backend/utils" // مهم جدًا عشان AutoMigrateAll
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// تحميل الـ .env لو موجود
	godotenv.Load()

	// الاتصال بقاعدة البيانات
	database.Connect()

	// إنشاء كل الجداول المتاحة (حاليًا بس users)
	log.Println("جاري إنشاء الجداول...")
	utils.AutoMigrateAll(database.DB)
	log.Println("تم إعداد قاعدة البيانات بنجاح")

	// إعداد Gin
	r := gin.Default()

	// تفعيل CORS لو هتشتغل من الموبايل أو الفرونت
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// الروتات
	api := r.Group("/api/v1")
	{
		routes.AuthRoutes(api)
		// لما تضيف موديلات وروتات جديدة هتضيفها هنا:
		// routes.ProductRoutes(api)
		// routes.OrderRoutes(api)
		// routes.CategoryRoutes(api)
	}

	// تشغيل السيرفر
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Eetmad Backend شغال 100%% على http://localhost:%s", port)
	log.Printf("جرب التسجيل: http://localhost:%s/api/v1/auth/register", port)
	log.Fatal(r.Run(":" + port))
}
