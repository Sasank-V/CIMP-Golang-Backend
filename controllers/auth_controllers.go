package controllers

import (
	"fmt"
	"os"
	"time"

	"github.com/Sasank-V/CIMP-Golang-Backend/types"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(payload types.TokenPayload) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      payload.ID,
		"name":    payload.Name,
		"is_lead": payload.IsLead,
		"exp":     time.Now().Add(time.Hour * 6).Unix(), // Expires in 6 hours
	})
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET_KEY is not set")
	}
	return token.SignedString([]byte(secret))
}

func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET_KEY is not set")
	}
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid Token")
}
