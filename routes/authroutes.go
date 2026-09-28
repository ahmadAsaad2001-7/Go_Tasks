package routes

import (
	"net/http"
	"task-manager-api/models"
	"task-manager-api/utils"

	"github.com/gin-gonic/gin"
)

func signup(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Default role to "user" if not provided
	if user.Role == "" {
		user.Role = "user"
	}

	createdUser, err := user.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

func login(c *gin.Context) {
	var userInput models.User

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Validate credentials AND get the user back in one query
	user, err := userInput.ValidateCredentials()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// 2. Generate Token directly using the returned user data!
	token, err := utils.GenerateToken(user.Email, user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
