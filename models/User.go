package models

import (
	"errors"
	"task-manager-api/db"
	"task-manager-api/utils"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email" binding:"required"`
	Password  string    `json:"password" binding:"required"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

func (u User) Save() (User, error) {
	query := `
	INSERT INTO users (email, password, role)
	VALUES ($1, $2, $3)
	RETURNING id, email, role, created_at
	`
	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return User{}, err
	}

	var user User
	err = db.DB.QueryRow(query, u.Email, hashedPassword, u.Role).Scan(
		&user.ID, &user.Email, &user.Role, &user.CreatedAt,
	)

	if err != nil {
		return User{}, err
	}
	return user, nil
}

// FIXED: Now returns (User, error) instead of just error
func (u User) ValidateCredentials() (User, error) {
	query := `
	SELECT id, email, password, role
	FROM users
	WHERE email = $1
	`
	var user User
	err := db.DB.QueryRow(query, u.Email).Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return User{}, err
	}

	isValidPassword := utils.ComparePasswords(user.Password, u.Password)
	if !isValidPassword {
		return User{}, errors.New("invalid credentials")
	}

	// Return the populated user struct!
	return user, nil
}

func (u User) FindByEmail() (User, error) {
	query := `
	SELECT id, email, role
	FROM users
	WHERE email = $1
	`
	var user User
	err := db.DB.QueryRow(query, u.Email).Scan(&user.ID, &user.Email, &user.Role)
	if err != nil {
		return User{}, err
	}
	return user, nil
}
