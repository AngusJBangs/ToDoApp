package main

import (
	"sync"

	"github.com/gin-gonic/gin"
)

func initialiseServer() Server {
	return Server{&sync.RWMutex{}, MakeSampleToDos(), DecodeToDo}
}

func main() {
	server := initialiseServer()
	router := gin.Default()
	router.GET("/todos", server.readAll)
	router.GET("/todo/:title", server.read)
	router.POST("/create", server.create)
	router.POST("/update/:title", server.update)
	router.POST("/delete/:title", server.delete)

	router.Run("localhost:8080")
}
