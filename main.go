package main

import (
	"ToDoApp/handlers"
	"ToDoApp/utilities"

	"github.com/gin-gonic/gin"
)

func initialiseServer() handlers.Server {
	return handlers.Server{StoredToDos: utilities.MakeSampleToDos(), DecodeToDo: utilities.DecodeToDo}
}

func main() {
	server := initialiseServer()
	server.Cmds = server.InitiateToDoHandlerManager()
	router := gin.Default()
	router.GET("/todos", server.StartReadAll)
	router.GET("/todo/:title", server.StartRead)
	router.POST("/create", server.StartCreate)
	router.POST("/update/:title", server.StartUpdate)
	router.POST("/delete/:title", server.StartDelete)

	router.Run("localhost:8080")
}
