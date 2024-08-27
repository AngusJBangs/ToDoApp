package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var l sync.Mutex
var sampleToDos = MakeSampleToDos()

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

func MakeSampleToDos() []ToDo {
	l.Lock()
	defer l.Unlock()
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

// func Jsonify(TD ...ToDo) ([]byte, bool) {
// 	jsonData, err := json.Marshal(TD)
// 	if err != nil {
// 		fmt.Println("Error marshaling to JSON:", err)
// 		return nil, false
// 	}
// 	return jsonData, true
// }

func readAll(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, sampleToDos)
}
func read(c *gin.Context) {
	title := c.Param("title")
	for _, td := range sampleToDos {
		if strings.EqualFold(strings.ReplaceAll(td.Title, " ", ""), title) {
			c.IndentedJSON(http.StatusOK, td)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "to-do not found"})
}
func create(c *gin.Context) {
	var tdr ToDoRequest
	if err := c.BindJSON(&tdr); err != nil {
		fmt.Println("Did not bind")
		return
	}
	due := strings.Split(tdr.Due, " ")
	year, err1 := strconv.Atoi(due[0])
	month, err2 := strconv.Atoi(due[1])
	day, err3 := strconv.Atoi(due[2])
	hour, err4 := strconv.Atoi(due[3])
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		fmt.Println("was not int")
		return
	}
	td := ToDo{tdr.Title, tdr.Description, time.Date(year, time.Month(month), day, hour, 0, 0, 0, time.Local), tdr.Priority, tdr.Status}
	l.Lock()
	defer l.Unlock()
	sampleToDos = append(sampleToDos, td)
	c.IndentedJSON(http.StatusOK, sampleToDos)
}
func update(c *gin.Context) {
	title := c.Param("title")
	var tdr ToDoRequest
	if err := c.BindJSON(&tdr); err != nil {
		fmt.Println("Did not bind")
		return
	}
	due := strings.Split(tdr.Due, " ")
	year, err1 := strconv.Atoi(due[0])
	month, err2 := strconv.Atoi(due[1])
	day, err3 := strconv.Atoi(due[2])
	hour, err4 := strconv.Atoi(due[3])
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		fmt.Println("was not int")
		return
	}
	newToDo := ToDo{tdr.Title, tdr.Description, time.Date(year, time.Month(month), day, hour, 0, 0, 0, time.Local), tdr.Priority, tdr.Status}
	l.Lock()
	defer l.Unlock()
	for i, td := range sampleToDos {
		if strings.EqualFold(strings.ReplaceAll(td.Title, " ", ""), title) {
			sampleToDos[i] = newToDo
			c.IndentedJSON(http.StatusOK, newToDo)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "to-do not found"})
}
func delete(c *gin.Context) {
	title := c.Param("title")
	l.Lock()
	defer l.Unlock()
	tds := append([]ToDo{}, sampleToDos...)
	for i, td := range sampleToDos {
		if strings.EqualFold(strings.ReplaceAll(td.Title, " ", ""), title) {
			before := tds[:i]
			after := []ToDo{}
			if i < len(tds) {
				after = tds[i+1:]
			}
			sampleToDos = append(before, after...)
			c.IndentedJSON(http.StatusOK, gin.H{"message": "to-do deleted"})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "to-do not found"})
}

func main() {
	router := gin.Default()
	router.GET("/todos", readAll)
	router.GET("/todo/:title", read)
	router.POST("/create", create)
	router.POST("/update/:title", update)
	router.POST("/delete/:title", delete)

	router.Run("localhost:8080")
}
