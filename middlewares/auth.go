package middlewares

import (
	"net/http"
	"task-manager-api/utils"

	"github.com/gin-gonic/gin"
)

type User struct {
	UserID int
	Email  string
	Role   string
}

func Authorization(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userInfo, err := utils.VerifyToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user := User{
		UserID: userInfo.UserID,
		Email:  userInfo.Email,
		Role:   userInfo.Role,
	}

	c.Set("user", user)
	c.Next()
} // <--- FIXED: Added this closing brace!

func AdminOnly(c *gin.Context) {
	userValue, exists := c.Get("user")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Type assertion to our User struct
	user, ok := userValue.(User)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type"})
		return
	}

	if user.Role != "admin" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	c.Next()
}
