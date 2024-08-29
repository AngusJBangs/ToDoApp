package models

import (
	"time"
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
	Priority    int    `json:"priority"`
	Status      string `json:"status"`
}
