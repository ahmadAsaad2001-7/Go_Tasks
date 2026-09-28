package routes

import (
	"net/http"
	"strconv"
	"task-manager-api/middlewares"
	"task-manager-api/models"

	"github.com/gin-gonic/gin"
)

func createTask(c *gin.Context) {
	projectID := c.Param("projectId")
	user := c.MustGet("user").(middlewares.User)

	// Verify user owns the project before adding a task
	project, err := models.GetProject(projectID)
	if err != nil || project.OwnerID != user.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You cannot add tasks to this project"})
		return
	}

	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	projIDInt, _ := strconv.Atoi(projectID)
	task.ProjectID = projIDInt

	createdTask, err := task.AddTask()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdTask)
}

// GET /projects/:projectId/tasks?status=todo&page=1&limit=10
func getTasks(c *gin.Context) {
	projectID := c.Param("projectId")
	statusFilter := c.Query("status")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	offset := (page - 1) * limit

	projIDInt, _ := strconv.Atoi(projectID)

	tasks, err := models.GetTasksByProject(projIDInt, statusFilter, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page":  page,
	"limit": limit,
	"data":  tasks,
	})
}

// ---------------------------------------------------------
// THE MISSING FUNCTIONS CAUSING YOUR COMPILER ERROR
// ---------------------------------------------------------

func updateTask(c *gin.Context) {
	id := c.Param("id")

	var taskInput models.Task
	if err := c.ShouldBindJSON(&taskInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch the existing task from the DB
	existingTask, err := models.GetTask(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Update the fields with the new input
	existingTask.Title = taskInput.Title
	existingTask.Description = taskInput.Description
	existingTask.Status = taskInput.Status
	existingTask.AssigneeID = taskInput.AssigneeID

	// Save to database
	updatedTask, err := existingTask.UpdateTask()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedTask)
}

func deleteTask(c *gin.Context) {
	id := c.Param("id")

	// Fetch the existing task first
	task, err := models.GetTask(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Soft delete it
	err = task.SoftDeleteTask()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
