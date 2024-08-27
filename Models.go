package main

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type ToDo struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Due         time.Time `json:"due"`
	Priority    int       `json:"priortiy"`
	Status      string    `json:"status"`
}

type ToDoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Due         string `json:"due"`
	Priority    int    `json:"priortiy"`
	Status      string `json:"status"`
}

type Server struct {
	l           *sync.RWMutex
	storedToDos []ToDo
	DecodeToDo  func(*gin.Context) (ToDo, error)
}
