package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func Test_readAll_HappyPath(t *testing.T) {

	initialToDos := []ToDo{
		{"Clean room", "Do laundry, hoover, clean sheets, change bin", time.Date(2024, 9, 15, 12, 0, 0, 0, time.Local), 2, "not started"},
		{"Clean car", "Remove mess, hoover, wash outside", time.Date(2024, 9, 12, 12, 0, 0, 0, time.Local), 4, "not started"}}
	server := Server{&sync.RWMutex{}, initialToDos, StubDecodeToDo(ToDo{}, nil)}
	response := httptest.NewRecorder()
	c := MockCustomContext("GET", ToDoRequest{}, []gin.Param{}, response)

	server.readAll(c)

	if response.Code != 200 {
		t.Error("Expected response to be 200 but was actually", response.Code)
	}
	var body []ToDo
	err := json.Unmarshal(response.Body.Bytes(), &body)
	if err != nil {
		t.Error("Expected response to be valid Json but could not unmarshall")
	}
	if b := body[0]; b != initialToDos[0] {
		t.Error("Expected", b, "to equal", initialToDos[0])
	}
	if b := body[1]; b != initialToDos[1] {
		t.Error("Expected", b, "to equal", initialToDos[1])
	}
	if l := len(body); l != 2 {
		t.Error("Expected response to have length 2 but actually of length", l)
	}
}
func Test_read_HappyPath(t *testing.T) {

	initialToDos := []ToDo{
		{"Clean room", "Do laundry, hoover, clean sheets, change bin", time.Date(2024, 9, 15, 12, 0, 0, 0, time.Local), 2, "not started"},
		{"Clean car", "Remove mess, hoover, wash outside", time.Date(2024, 9, 12, 12, 0, 0, 0, time.Local), 4, "not started"}}
	server := Server{&sync.RWMutex{}, initialToDos, StubDecodeToDo(ToDo{}, nil)}
	response := httptest.NewRecorder()
	c := MockCustomContext("GET", ToDoRequest{}, []gin.Param{{Key: "title", Value: "cleanroom"}}, response)

	server.read(c)

	if response.Code != 200 {
		t.Error("Expected response to be 200 but was actually", response.Code)
	}
	var body ToDo
	err := json.Unmarshal(response.Body.Bytes(), &body)
	if err != nil {
		t.Error("Expected response to be valid Json but could not unmarshall")
	}
	if b := body; b != initialToDos[0] {
		t.Error("Expected", b, "to equal", initialToDos[0])
	}
}
func Test_read_404TitleNotFound(t *testing.T) {

	initialToDos := []ToDo{
		{"Clean room", "Do laundry, hoover, clean sheets, change bin", time.Date(2024, 9, 15, 12, 0, 0, 0, time.Local), 2, "not started"},
		{"Clean car", "Remove mess, hoover, wash outside", time.Date(2024, 9, 12, 12, 0, 0, 0, time.Local), 4, "not started"}}
	server := Server{&sync.RWMutex{}, initialToDos, StubDecodeToDo(ToDo{}, nil)}
	response := httptest.NewRecorder()
	c := MockCustomContext("GET", ToDoRequest{}, []gin.Param{{Key: "title", Value: "slackOff"}}, response)

	server.read(c)

	if response.Code != 404 {
		t.Error("Expected response to be 404 but was actually", response.Code)
	}
}
func Test_delete_HappyPath(t *testing.T) {

	initialToDos := []ToDo{
		{"Clean room", "Do laundry, hoover, clean sheets, change bin", time.Date(2024, 9, 15, 12, 0, 0, 0, time.Local), 2, "not started"},
		{"Clean car", "Remove mess, hoover, wash outside", time.Date(2024, 9, 12, 12, 0, 0, 0, time.Local), 4, "not started"}}
	server := Server{&sync.RWMutex{}, initialToDos, StubDecodeToDo(ToDo{}, nil)}
	response := httptest.NewRecorder()
	c := MockCustomContext("POST", ToDoRequest{}, []gin.Param{{Key: "title", Value: "CleanCar"}}, response)

	server.delete(c)

	if response.Code != 200 {
		t.Error("Expected response to be 200 but was actually", response.Code)
	}
	if l := len(server.storedToDos); l != 1 {
		t.Error("Expected only 1 to-do to remain after deletion but was ", l)
	}
	if s := server.storedToDos[0]; s != initialToDos[0] {
		t.Error("Expected remaining to-do", s, "to equal", initialToDos[0])
	}
}
func Test_delete_404TitleNotFound(t *testing.T) {

	initialToDos := []ToDo{
		{"Clean room", "Do laundry, hoover, clean sheets, change bin", time.Date(2024, 9, 15, 12, 0, 0, 0, time.Local), 2, "not started"},
		{"Clean car", "Remove mess, hoover, wash outside", time.Date(2024, 9, 12, 12, 0, 0, 0, time.Local), 4, "not started"}}
	server := Server{&sync.RWMutex{}, initialToDos, StubDecodeToDo(ToDo{}, nil)}
	response := httptest.NewRecorder()
	c := MockCustomContext("POST", ToDoRequest{}, []gin.Param{{Key: "title", Value: "takeItEasy"}}, response)

	server.delete(c)

	if response.Code != 404 {
		t.Error("Expected response to be 404 but was actually", response.Code)
	}
}
func Test_update_HappyPath(t *testing.T) {

	initialToDos := []ToDo{
		{"Clean room", "Do laundry, hoover, clean sheets, change bin", time.Date(2024, 9, 15, 12, 0, 0, 0, time.Local), 2, "not started"},
		{"Clean car", "Remove mess, hoover, wash outside", time.Date(2024, 9, 12, 12, 0, 0, 0, time.Local), 4, "not started"}}
	expected := ToDo{"Clean room", "an updated description", time.Date(2024, 9, 15, 12, 0, 0, 0, time.Local), 2, "not started"}
	server := Server{&sync.RWMutex{}, initialToDos, StubDecodeToDo(expected, nil)}
	response := httptest.NewRecorder()
	c := MockCustomContext("POST", ToDoRequest{"Clean room", "an updated description", "2024 9 15 12", 2, "not started"}, []gin.Param{{Key: "title", Value: "Clean Room"}}, response)

	server.update(c)

	if response.Code != 200 {
		t.Error("Expected response to be 200 but was actually", response.Code)
	}
	if l := len(server.storedToDos); l != 2 {
		t.Error("Expected 2 to-dos to remain after update but was ", l)
	}
	if s := server.storedToDos[0]; s != expected {
		t.Error("Expected remaining to-do", s, "to equal", expected)
	}
	if s := server.storedToDos[1]; s != initialToDos[1] {
		t.Error("Expected remaining to-do", s, "to equal", initialToDos[1])
	}
}
func Test_update_404TitleNotFound(t *testing.T) {

	initialToDos := []ToDo{
		{"Clean room", "Do laundry, hoover, clean sheets, change bin", time.Date(2024, 9, 15, 12, 0, 0, 0, time.Local), 2, "not started"},
		{"Clean car", "Remove mess, hoover, wash outside", time.Date(2024, 9, 12, 12, 0, 0, 0, time.Local), 4, "not started"}}
	expected := ToDo{"Clean room", "an updated description", time.Date(2024, 9, 15, 12, 0, 0, 0, time.Local), 2, "not started"}
	server := Server{&sync.RWMutex{}, initialToDos, StubDecodeToDo(expected, nil)}
	response := httptest.NewRecorder()
	c := MockCustomContext("POST", ToDoRequest{"Clean room", "an updated description", "2024 9 15 12", 2, "not started"}, []gin.Param{{Key: "title", Value: "SnoozeAlarm"}}, response)

	server.update(c)

	if response.Code != 404 {
		t.Error("Expected response to be 404 but was actually", response.Code)
	}
}
func Test_update_409TitleAlreadyExists(t *testing.T) {

	initialToDos := []ToDo{
		{"Clean room", "Do laundry, hoover, clean sheets, change bin", time.Date(2024, 9, 15, 12, 0, 0, 0, time.Local), 2, "not started"},
		{"Clean car", "Remove mess, hoover, wash outside", time.Date(2024, 9, 12, 12, 0, 0, 0, time.Local), 4, "not started"}}
	expected := ToDo{"Clean car", "an updated description", time.Date(2024, 9, 15, 12, 0, 0, 0, time.Local), 2, "not started"}
	server := Server{&sync.RWMutex{}, initialToDos, StubDecodeToDo(expected, nil)}
	response := httptest.NewRecorder()
	c := MockCustomContext("POST", ToDoRequest{"Clean Car", "an updated description", "2024 9 15 12", 2, "not started"}, []gin.Param{{Key: "title", Value: "cleanRoom"}}, response)

	server.update(c)

	if response.Code != 409 {
		t.Error("Expected response to be 409 but was actually", response.Code)
	}
}
func Test_update_400BadRequest(t *testing.T) {

	initialToDos := []ToDo{
		{"Clean room", "Do laundry, hoover, clean sheets, change bin", time.Date(2024, 9, 15, 12, 0, 0, 0, time.Local), 2, "not started"},
		{"Clean car", "Remove mess, hoover, wash outside", time.Date(2024, 9, 12, 12, 0, 0, 0, time.Local), 4, "not started"}}
	server := Server{&sync.RWMutex{}, initialToDos, StubDecodeToDo(ToDo{}, errors.New("an error message"))}
	response := httptest.NewRecorder()
	c := MockCustomContext("POST", ToDoRequest{}, []gin.Param{}, response)

	server.update(c)

	if response.Code != 400 {
		t.Error("Expected response to be 400 but was actually", response.Code)
	}
	if m := strings.Split(response.Body.String(), `"`)[3]; m != "an error message" {
		t.Error(`Expected response message to be "an error message" but was actually`, m)
	}
}
func Test_create_HappyPath(t *testing.T) {

	initialToDos := []ToDo{
		{"Clean room", "Do laundry, hoover, clean sheets, change bin", time.Date(2024, 9, 15, 12, 0, 0, 0, time.Local), 2, "not started"},
		{"Clean car", "Remove mess, hoover, wash outside", time.Date(2024, 9, 12, 12, 0, 0, 0, time.Local), 4, "not started"}}
	expected := ToDo{"new title", "new description", time.Date(2024, 1, 1, 1, 0, 0, 0, time.Local), 1, "not started"}
	server := Server{&sync.RWMutex{}, initialToDos, StubDecodeToDo(expected, nil)}
	response := httptest.NewRecorder()
	c := MockCustomContext("POST", ToDoRequest{"new title", "new description", "2024 1 1 1", 1, "not started"}, []gin.Param{}, response)

	server.create(c)

	if response.Code != 200 {
		t.Error("Expected response to be 200 but was actually", response.Code)
	}
	if l := len(server.storedToDos); l != 3 {
		t.Error("Expected 3 to-dos to exist after creation but was ", l)
	}
	if s := server.storedToDos[2]; s != expected {
		t.Error("Expected new to-do", s, "to equal", expected)
	}
	if s := server.storedToDos[1]; s != initialToDos[1] {
		t.Error("Expected existing to-do", s, "to equal", initialToDos[1])
	}
	if s := server.storedToDos[0]; s != initialToDos[0] {
		t.Error("Expected existing to-do", s, "to equal", initialToDos[0])
	}
}
func Test_create_409TitleAlreadyExists(t *testing.T) {

	initialToDos := []ToDo{
		{"Clean room", "Do laundry, hoover, clean sheets, change bin", time.Date(2024, 9, 15, 12, 0, 0, 0, time.Local), 2, "not started"},
		{"Clean car", "Remove mess, hoover, wash outside", time.Date(2024, 9, 12, 12, 0, 0, 0, time.Local), 4, "not started"}}
	expected := ToDo{"clean room", "new description", time.Date(2024, 1, 1, 1, 0, 0, 0, time.Local), 1, "not started"}
	server := Server{&sync.RWMutex{}, initialToDos, StubDecodeToDo(expected, nil)}
	response := httptest.NewRecorder()
	c := MockCustomContext("POST", ToDoRequest{"clean room", "new description", "2024 1 1 1", 1, "not started"}, []gin.Param{}, response)

	server.create(c)

	if response.Code != 409 {
		t.Error("Expected response to be 409 but was actually", response.Code)
	}
}
func Test_create_400BadRequest(t *testing.T) {

	initialToDos := []ToDo{
		{"Clean room", "Do laundry, hoover, clean sheets, change bin", time.Date(2024, 9, 15, 12, 0, 0, 0, time.Local), 2, "not started"},
		{"Clean car", "Remove mess, hoover, wash outside", time.Date(2024, 9, 12, 12, 0, 0, 0, time.Local), 4, "not started"}}
	server := Server{&sync.RWMutex{}, initialToDos, StubDecodeToDo(ToDo{}, errors.New("missing data"))}
	response := httptest.NewRecorder()
	c := MockCustomContext("POST", ToDoRequest{}, []gin.Param{}, response)

	server.create(c)

	if response.Code != 400 {
		t.Error("Expected response to be 400 but was actually", response.Code)
	}
	if m := strings.Split(response.Body.String(), `"`)[3]; m != "missing data" {
		t.Error(`Expected response message to be "missing data" but was actually`, m)
	}
}

func MockCustomContext(methodType string, tdr ToDoRequest, p []gin.Param, w *httptest.ResponseRecorder) *gin.Context {
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
	}
	c.Request.Method = methodType
	c.Params = p
	if methodType == "POST" {
		c.Request.Header.Set("Content-Type", "application/json")
		jsonbytes, err := json.Marshal(tdr)
		if err != nil {
			panic(err)
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
	}
	return c
}

func StubDecodeToDo(td ToDo, err error) func(c *gin.Context) (ToDo, error) {
	return func(c *gin.Context) (ToDo, error) {
		return td, err
	}
}
