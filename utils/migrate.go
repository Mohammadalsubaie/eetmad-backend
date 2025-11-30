// utils/migrate.go
package utils

import (
	"log"

	"github.com/eetmad/backend/models"
	"gorm.io/gorm"
)

// AutoMigrateAll → تعمل كل الجداول الموجودة فعليًا فقط
func AutoMigrateAll(db *gorm.DB) {
	log.Println("بدء إنشاء الجداول المتاحة...")

	// هنعمل Migrate لكل الموديلات اللي موجودة فعليًا عندك
	// حاليًا بس الـ User موجود → وده كفاية جدًا للـ MVP
	db.AutoMigrate(&models.User{})

	// لو في المستقبل ضفت موديلات جديدة (Product, Order, ...) ضيفها هنا
	// مثال:
	// db.AutoMigrate(&models.Product{}, &models.Category{}, ...)

	log.Println("تم إنشاء جميع الجداول المتاحة بنجاح")
}
