package main

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func MakeSampleToDos() []ToDo {
	return []ToDo{
		{"Clean room", "Do laundry, hoover, clean sheets, change bin", time.Date(2024, 9, 15, 12, 0, 0, 0, time.Local), 2, "not started"},
		{"Clean car", "Remove mess, hoover, wash outside", time.Date(2024, 9, 12, 12, 0, 0, 0, time.Local), 4, "not started"},
		{"Buy Present Mum", "Buy a birthday present for mum, maybe a massage(?)", time.Date(2024, 9, 25, 12, 0, 0, 0, time.Local), 9, "not started"},
		{"ChaseBuilders", "Chase for quote builder@builder.com", time.Date(2024, 9, 11, 12, 0, 0, 0, time.Local), 3, "in progress"},
		{"Finish Go Academy", "Complete ToDoApp inc testing", time.Date(2024, 8, 30, 12, 30, 0, 0, time.Local), 10, "not started"},
		{"Learn to fly", "Jump off progressively higher things until I figure it out", time.Date(2024, 9, 28, 12, 0, 0, 0, time.Local), 3, "not started"},
		{"Swallow spider", "Just something big enough to catch that fly", time.Date(2024, 8, 29, 12, 0, 0, 0, time.Local), 4, "complete"},
		{"Nigerian Prince", "Nigerian prince is driving around with a truck full of gold, waiting for me to send money sam.cam@hotmail.com", time.Date(2024, 9, 1, 12, 0, 0, 0, time.Local), 7, "not started"},
		{"Count to a billion", "Start at 1, count to a billion, tell theboy next door that I counted to a bigger number than him", time.Date(2024, 9, 9, 12, 0, 0, 0, time.Local), 10, "not started"},
		{"Leave the house", "Probably a bit of a stretch goal but worth a try", time.Date(2050, 1, 1, 00, 0, 0, 0, time.Local), 1, "not started"},
	}
}

func DecodeToDo(c *gin.Context) (ToDo, error) {
	var tdr ToDoRequest
	if err := c.BindJSON(&tdr); err != nil {
		return ToDo{}, errors.New("to-do incorrectly formatted")
	}
	if tdr.Title == "" || tdr.Description == "" || tdr.Due == "" || tdr.Priority == 0 || tdr.Status == "" {
		return ToDo{}, errors.New("missing data")
	}
	due := strings.Split(tdr.Due, " ")
	if len(due) < 4 {
		return ToDo{}, errors.New("due incorrectly formatted - not enough vales for date")
	}
	year, err1 := strconv.Atoi(due[0])
	month, err2 := strconv.Atoi(due[1])
	day, err3 := strconv.Atoi(due[2])
	hour, err4 := strconv.Atoi(due[3])
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return ToDo{}, errors.New("due incorrectly formatted - value not int")
	}
	if month < 1 || day < 1 || hour < 1 || month > 12 || day > 31 || hour > 24 {
		return ToDo{}, errors.New("due incorrectly formatted - value out of range")
	}

	return ToDo{tdr.Title, tdr.Description, time.Date(year, time.Month(month), day, hour, 0, 0, 0, time.Local), tdr.Priority, tdr.Status}, nil
}
