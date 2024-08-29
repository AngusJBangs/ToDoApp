package handlers

import (
	"ToDoApp/models"
	"sync"

	"github.com/gin-gonic/gin"
)

type Server struct {
	L           *sync.RWMutex
	StoredToDos []models.ToDo
	DecodeToDo  func(*gin.Context) (models.ToDo, error)
}
