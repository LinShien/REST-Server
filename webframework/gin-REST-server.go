package main

import (
	"github.com/gin-gonic/gin"
	"github.com/shien/restserver/webframework/taskserver"
)

func main() {
	router := gin.Default()
	server := taskserver.NewTaskServerForWebFramework()

	// register, unlike Router package, there is no regexp support in Gin(Web framework)
	router.GET("/task/", server.GetAllTasksHandler)
	router.POST("/task/", server.CreateTaskHandler)
	router.DELETE("/task/", server.DeleteAllTasksHandler)

	router.GET("/task/:id", server.GetTaskHandler)
	router.DELETE("/task/:id", server.DeleteTaskHandler)

	router.GET("/tag/:tag", server.TagHandler)
	router.GET("/due/:year/:month/:day", server.DueHandler)

	const PORT = "9090"

	router.Run("localhost:" + PORT)
}
