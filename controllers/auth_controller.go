package controllers

import (
	"net/http"
	"regexp"
	"github.com/eetmad/backend/database"
	"github.com/eetmad/backend/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=3,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone,omitempty"`
	Password string `json:"password" binding:"required,min=8,max=128"`
	UserType string `json:"user_type,omitempty"`
}

type RegisterResponse struct {
	Message string `json:"message"`
	User    UserData `json:"user"`
	Token   string `json:"token,omitempty"`
}

type UserData struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	UserType string `json:"user_type"`
}

func Register(c *gin.Context) {
	var input RegisterRequest

	// Bind JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate email format
	if matched, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, input.Email); !matched {
		c.JSON(http.StatusBadRequest, gin.H{"error": "البريد الإلكتروني غير صالح"})
		return
	}

	// Validate password strength (optional but recommended)
	if len(input.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "كلمة المرور يجب أن تكون 8 أحرف على الأقل"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطأ في معالجة كلمة المرور"})
		return
	}

	// Determine user type
	userType := "client"
	if input.UserType == "supplier" {
		userType = "supplier"
	}

	// Create user
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Phone:    input.Phone,
		Password: string(hashedPassword),
		UserType: userType,
	}

	// Save to database
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "البريد مستخدم مسبقًا"})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, RegisterResponse{
		Message: "تم التسجيل بنجاح",
		User: UserData{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			UserType: user.UserType,
		},
	})
}
