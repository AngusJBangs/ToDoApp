package main

import (
	"ToDoApp/models"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func Test_decodeToDo_HappyPath(t *testing.T) {

	c := MockPostContext(models.ToDoRequest{Title: "title", Description: "description", Due: "2024 9 15 12", Priority: 5, Status: "status"})

	result, err := DecodeToDo(c)

	if err != nil {
		t.Error("expected happy path but errored with ", err.Error())
	}
	expected := models.ToDo{Title: "title", Description: "description", Due: time.Date(2024, 9, 15, 12, 0, 0, 0, time.Local), Priority: 5, Status: "status"}
	if result != expected {
		t.Error("Expected to-do to be ", expected, " but was actuall ", result)
	}
}

func Test_decodeToDo_ErrorsWhenNegativeDate(t *testing.T) {
	c := MockPostContext(models.ToDoRequest{Title: "title", Description: "description", Due: "2024 9 -15 12", Priority: 5, Status: "status"})
	_, err := DecodeToDo(c)
	AssertError(err, "due incorrectly formatted - value out of range", t)
}
func Test_decodeToDo_ErrorsWhenMonthSetTo13(t *testing.T) {
	c := MockPostContext(models.ToDoRequest{Title: "title", Description: "description", Due: "2024 -9 15 12", Priority: 5, Status: "status"})
	_, err := DecodeToDo(c)

	AssertError(err, "due incorrectly formatted - value out of range", t)
}
func Test_decodeToDo_ErrorsWhenIncompleteDueDate(t *testing.T) {
	c := MockPostContext(models.ToDoRequest{Title: "title", Description: "description", Due: "2024 9 15", Priority: 5, Status: "status"})
	_, err := DecodeToDo(c)
	AssertError(err, "due incorrectly formatted - not enough vales for date", t)
}
func Test_decodeToDo_ErrorsWhenDateNotInt(t *testing.T) {
	c := MockPostContext(models.ToDoRequest{Title: "title", Description: "description", Due: "2024 september 15 12", Priority: 5, Status: "status"})
	_, err := DecodeToDo(c)
	AssertError(err, "due incorrectly formatted - value not int", t)
}
func Test_decodeToDo_ErrorsWhenInputInvalid(t *testing.T) {

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.Method = "POST"
	c.Request.Header.Set("Content-Type", "application/json")
	jsonbytes, err := json.Marshal(incompleteToDo{"lots of missing values"})
	if err != nil {
		panic(err)
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
	_, err = DecodeToDo(c)
	AssertError(err, "missing data", t)
}

func FuzzDecodeToDo(f *testing.F) {
	f.Add("title", "description", "2024 8 6 4", 4, "not started")
	f.Fuzz(func(t *testing.T, title string, description string, due string, priority int, status string) {
		c := MockPostContext(models.ToDoRequest{Title: title, Description: description, Due: due, Priority: priority, Status: status})
		_, err := DecodeToDo(c)
		if err != nil {
			t.Skip("invalid entry")
		}
		t.Skip("valid entry")
	})
}

func AssertError(err error, expectedErrorMessage string, t *testing.T) {
	if err == nil {
		t.Error("Expected error with message:'", expectedErrorMessage, "' but no error recieved")
		return
	}
	if err.Error() != expectedErrorMessage {
		t.Error("Expected error message to be:'", expectedErrorMessage, "' but was '", err.Error(), "'")
		return
	}
}

func MockPostContext(tdr models.ToDoRequest) *gin.Context {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.Method = "POST"
	c.Request.Header.Set("Content-Type", "application/json")
	jsonbytes, err := json.Marshal(tdr)
	if err != nil {
		panic(err)
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
	return c
}

type incompleteToDo struct {
	Something string `json:"something"`
}
