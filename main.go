package main

import (
	"ToDoApp/handlers"
	"sync"

	"github.com/gin-gonic/gin"
)

func initialiseServer() handlers.Server {
	return handlers.Server{L: &sync.RWMutex{}, StoredToDos: MakeSampleToDos(), DecodeToDo: DecodeToDo}
}

func main() {
	server := initialiseServer()
	router := gin.Default()
	router.GET("/todos", server.ReadAll)
	router.GET("/todo/:title", server.Read)
	router.POST("/create", server.Create)
	router.POST("/update/:title", server.Update)
	router.POST("/delete/:title", server.Delete)

	router.Run("localhost:8080")
}
