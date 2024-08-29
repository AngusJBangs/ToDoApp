package main

import (
	"ToDoApp/models"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func MakeSampleToDos() []models.ToDo {
	return []models.ToDo{
		{Title: "Clean room", Description: "Do laundry, hoover, clean sheets, change bin", Due: time.Date(2024, 9, 15, 12, 0, 0, 0, time.Local), Priority: 2, Status: "not started"},
		{Title: "Clean car", Description: "Remove mess, hoover, wash outside", Due: time.Date(2024, 9, 12, 12, 0, 0, 0, time.Local), Priority: 4, Status: "not started"},
		{Title: "Buy Present Mum", Description: "Buy a birthday present for mum, maybe a massage(?)", Due: time.Date(2024, 9, 25, 12, 0, 0, 0, time.Local), Priority: 9, Status: "not started"},
		{Title: "ChaseBuilders", Description: "Chase for quote builder@builder.com", Due: time.Date(2024, 9, 11, 12, 0, 0, 0, time.Local), Priority: 3, Status: "in progress"},
		{Title: "Finish Go Academy", Description: "Complete ToDoApp inc testing", Due: time.Date(2024, 8, 30, 12, 30, 0, 0, time.Local), Priority: 10, Status: "not started"},
		{Title: "Learn to fly", Description: "Jump off progressively higher things until I figure it out", Due: time.Date(2024, 9, 28, 12, 0, 0, 0, time.Local), Priority: 3, Status: "not started"},
		{Title: "Swallow spider", Description: "Just something big enough to catch that fly", Due: time.Date(2024, 8, 29, 12, 0, 0, 0, time.Local), Priority: 4, Status: "complete"},
		{Title: "Nigerian Prince", Description: "Nigerian prince is driving around with a truck full of gold, waiting for me to send money sam.cam@hotmail.com", Due: time.Date(2024, 9, 1, 12, 0, 0, 0, time.Local), Priority: 7, Status: "not started"},
		{Title: "Count to a billion", Description: "Start at 1, count to a billion, tell theboy next door that I counted to a bigger number than him", Due: time.Date(2024, 9, 9, 12, 0, 0, 0, time.Local), Priority: 10, Status: "not started"},
		{Title: "Leave the house", Description: "Probably a bit of a stretch goal but worth a try", Due: time.Date(2050, 1, 1, 00, 0, 0, 0, time.Local), Priority: 1, Status: "not started"},
	}
}

func DecodeToDo(c *gin.Context) (models.ToDo, error) {
	var tdr models.ToDoRequest
	if err := c.BindJSON(&tdr); err != nil {
		return models.ToDo{}, errors.New("to-do incorrectly formatted")
	}
	if tdr.Title == "" || tdr.Description == "" || tdr.Due == "" || tdr.Priority == 0 || tdr.Status == "" {
		return models.ToDo{}, errors.New("missing data")
	}
	due := strings.Split(tdr.Due, " ")
	if len(due) < 4 {
		return models.ToDo{}, errors.New("due incorrectly formatted - not enough vales for date")
	}
	year, err1 := strconv.Atoi(due[0])
	month, err2 := strconv.Atoi(due[1])
	day, err3 := strconv.Atoi(due[2])
	hour, err4 := strconv.Atoi(due[3])
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return models.ToDo{}, errors.New("due incorrectly formatted - value not int")
	}
	if month < 1 || day < 1 || hour < 1 || month > 12 || day > 31 || hour > 24 {
		return models.ToDo{}, errors.New("due incorrectly formatted - value out of range")
	}

	return models.ToDo{Title: tdr.Title, Description: tdr.Description, Due: time.Date(year, time.Month(month), day, hour, 0, 0, 0, time.Local), Priority: tdr.Priority, Status: tdr.Status}, nil
}
