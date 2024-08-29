package handlers

import (
	"ToDoApp/models"

	"github.com/gin-gonic/gin"
)

type Server struct {
	StoredToDos []models.ToDo
	DecodeToDo  func(*gin.Context) (models.ToDo, error)
	Cmds        chan<- Command
}
