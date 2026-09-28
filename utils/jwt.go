package utils

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken now accepts a role parameter
func GenerateToken(email string, userId int, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"role":   role, // Include the role in the token
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})

	JWT_SECRET := os.Getenv("JWT_SECRET")
	if JWT_SECRET == "" {
		return "", errors.New("JWT secret is not configured")
	}

	return token.SignedString([]byte(JWT_SECRET))
}

// VerifyToken now returns a struct that includes the Role
func VerifyToken(tokenString string) (struct {
	UserID int
	Email  string
	Role   string
}, error) {
	result := struct {
		UserID int
		Email  string
		Role   string
	}{}

	// Fixed typo: TrimPrefix
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	JWT_SECRET := os.Getenv("JWT_SECRET")
	if JWT_SECRET == "" {
		return result, errors.New("JWT secret is not configured")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Fixed type assertion syntax
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JWT_SECRET), nil
	})

	if err != nil {
		return result, fmt.Errorf("token parsing error: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return result, errors.New("invalid token claims")
	}

	userIDFloat, ok := claims["userId"].(float64)
	if !ok || userIDFloat == 0 {
		return result, errors.New("invalid or missing user ID")
	}
	// Fixed: Assign the UserID to the result
	result.UserID = int(userIDFloat)

	email, ok := claims["email"].(string)
	if !ok || email == "" {
		return result, errors.New("invalid or missing email")
	}
	result.Email = email

	// Extract the role
	role, ok := claims["role"].(string)
	if !ok {
		return result, errors.New("invalid or missing role")
	}
	result.Role = role

	return result, nil
}
