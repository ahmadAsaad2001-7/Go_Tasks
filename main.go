package main

import (
	"task-manager-api/db"
	"task-manager-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	server := gin.Default()

	routes.RegisterRoutes(server)

	server.Run(":4001")
}
