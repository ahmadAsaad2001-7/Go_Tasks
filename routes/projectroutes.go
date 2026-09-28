package routes

import (
	"net/http"
	"task-manager-api/middlewares"
	"task-manager-api/models"

	"github.com/gin-gonic/gin"
)

func createProject(c *gin.Context) {
	user := c.MustGet("user").(middlewares.User)

	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project.OwnerID = user.UserID

	createdProject, err := project.AddProject()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdProject)
}

func getMyProjects(c *gin.Context) {
	user := c.MustGet("user").(middlewares.User)

	projects, err := models.GetProjectsByUser(user.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projects)
}

func updateProject(c *gin.Context) {
	id := c.Param("id")
	user := c.MustGet("user").(middlewares.User)

	project, err := models.GetProject(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	if project.OwnerID != user.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not own this project"})
		return
	}

	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project.OwnerID = user.UserID
	updatedProject, err := project.UpdateProject()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedProject)
}

func deleteProject(c *gin.Context) {
	id := c.Param("id")
	user := c.MustGet("user").(middlewares.User)

	project, err := models.GetProject(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	if project.OwnerID != user.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not own this project"})
		return
	}

	err = project.SoftDeleteProject()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ---------------------------------------------------------
// THE MISSING FUNCTIONS THAT FIXED YOUR COMPILER ERROR
// ---------------------------------------------------------

func getProject(c *gin.Context) {
	id := c.Param("id")
	user := c.MustGet("user").(middlewares.User)

	project, err := models.GetProject(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Allow access if the user is the owner OR if they are an admin
	if project.OwnerID != user.UserID && user.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not own this project"})
		return
	}

	c.JSON(http.StatusOK, project)
}

func adminDeleteProject(c *gin.Context) {
	id := c.Param("id")

	project, err := models.GetProject(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Admin bypasses the ownership check. The AdminOnly middleware already verified they are an admin.
	err = project.SoftDeleteProject()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
