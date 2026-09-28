package routes

import (
	"task-manager-api/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	// Public Routes
	router.POST("/auth/signup", signup)
	router.POST("/auth/login", login)

	// Protected Routes (Must be logged in)
	authenticated := router.Group("/")
	authenticated.Use(middlewares.Authorization)

	// Project Routes
	authenticated.POST("/projects", createProject)
	authenticated.GET("/projects", getMyProjects)
	authenticated.GET("/projects/:id", getProject)
	authenticated.PUT("/projects/:id", updateProject)
	authenticated.DELETE("/projects/:id", deleteProject)

	// Task Routes
	authenticated.POST("/projects/:id/tasks", createTask)
	authenticated.GET("/projects/:id/tasks", getTasks)
	authenticated.PUT("/tasks/:id", updateTask)
	authenticated.DELETE("/tasks/:id", deleteTask)

	// Admin Only Route (Example)
	authenticated.DELETE("/admin/projects/:id", middlewares.AdminOnly, adminDeleteProject)
}
