package utils

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   string `json:"user_id"`
	UserType string `json:"user_type"`
	jwt.RegisteredClaims
}

func GenerateToken(userID interface{}, userType ...string) (string, error) {
	var idStr string
	switch v := userID.(type) {
	case uint:
		idStr = strconv.FormatUint(uint64(v), 10)
	case uint64:
		idStr = strconv.FormatUint(v, 10)
	case int:
		idStr = strconv.Itoa(v)
	case string:
		idStr = v
	default:
		idStr = "0"
	}

	ut := "client"
	if len(userType) > 0 {
		ut = userType[0]
	}

	claims := Claims{
		UserID:   idStr,
		UserType: ut,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "eetmad",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
