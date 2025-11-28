package controllers

import (
	"net/http"
	"os"
	"regexp"

	"github.com/eetmad/backend/database"
	"github.com/eetmad/backend/models"
	"github.com/eetmad/backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone,omitempty"`
	Password string `json:"password" binding:"required,min=8"`
	UserType string `json:"user_type,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// Register
func Register(c *gin.Context) {
	var input RegisterRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if matched, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, input.Email); !matched {
		c.JSON(http.StatusBadRequest, gin.H{"error": "البريد الإلكتروني غير صالح"})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	userType := "client"
	if input.UserType == "supplier" || input.UserType == "admin" {
		userType = input.UserType
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Phone:    input.Phone,
		Password: string(hashed),
		UserType: userType,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "البريد الإلكتروني مستخدم مسبقًا"})
		return
	}

	token, _ := utils.GenerateToken(user.ID, user.UserType)

	c.JSON(http.StatusCreated, gin.H{
		"message": "تم التسجيل بنجاح",
		"user":    gin.H{"id": user.ID, "name": user.Name, "email": user.Email, "user_type": user.UserType},
		"token":   token,
	})
}

// Login
func Login(c *gin.Context) {
	var input LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "بيانات غير صحيحة"})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "البريد أو كلمة المرور غير صحيحة"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "البريد أو كلمة المرور غير صحيحة"})
		return
	}

	token, _ := utils.GenerateToken(user.ID, user.UserType)
	refreshToken, _ := utils.GenerateToken(user.ID, user.UserType)

	c.SetCookie("refresh_token", refreshToken, 30*24*3600, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "تم تسجيل الدخول بنجاح",
		"token":   token,
	})
}

// Refresh Token
func Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "توكن منتهي"})
		return
	}

	claims := &utils.Claims{}
	_, err = jwt.ParseWithClaims(refreshToken, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "توكن غير صالح"})
		return
	}

	newToken, _ := utils.GenerateToken(claims.UserID, claims.UserType)
	c.JSON(http.StatusOK, gin.H{"token": newToken})
}

// Me
func Me(c *gin.Context) {
	userID := c.GetString("user_id")
	var user models.User
	database.DB.Select("id, name, email, phone, user_type, created_at").First(&user, userID)
	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"phone":      user.Phone,
		"user_type":  user.UserType,
		"created_at": user.CreatedAt,
	})
}
