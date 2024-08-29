package handlers

import (
	"ToDoApp/models"
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

func Test_ReadAll_HappyPath(t *testing.T) {

	initialToDos := []models.ToDo{
		{Title: "Clean room", Description: "Do laundry, hoover, clean sheets, change bin", Due: time.Date(2024, 9, 15, 12, 0, 0, 0, time.Local), Priority: 2, Status: "not started"},
		{Title: "Clean car", Description: "Remove mess, hoover, wash outside", Due: time.Date(2024, 9, 12, 12, 0, 0, 0, time.Local), Priority: 4, Status: "not started"}}
	server := Server{&sync.RWMutex{}, initialToDos, StubDecodeToDo(models.ToDo{}, nil)}
	response := httptest.NewRecorder()
	c := MockCustomContext("GET", models.ToDoRequest{}, []gin.Param{}, response)

	server.ReadAll(c)

	if response.Code != 200 {
		t.Error("Expected response to be 200 but was actually", response.Code)
	}
	var body []models.ToDo
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
func Test_Read_HappyPath(t *testing.T) {

	initialToDos := []models.ToDo{
		{Title: "Clean room", Description: "Do laundry, hoover, clean sheets, change bin", Due: time.Date(2024, 9, 15, 12, 0, 0, 0, time.Local), Priority: 2, Status: "not started"},
		{Title: "Clean car", Description: "Remove mess, hoover, wash outside", Due: time.Date(2024, 9, 12, 12, 0, 0, 0, time.Local), Priority: 4, Status: "not started"}}
	server := Server{&sync.RWMutex{}, initialToDos, StubDecodeToDo(models.ToDo{}, nil)}
	response := httptest.NewRecorder()
	c := MockCustomContext("GET", models.ToDoRequest{}, []gin.Param{{Key: "title", Value: "cleanroom"}}, response)

	server.Read(c)

	if response.Code != 200 {
		t.Error("Expected response to be 200 but was actually", response.Code)
	}
	var body models.ToDo
	err := json.Unmarshal(response.Body.Bytes(), &body)
	if err != nil {
		t.Error("Expected response to be valid Json but could not unmarshall")
	}
	if b := body; b != initialToDos[0] {
		t.Error("Expected", b, "to equal", initialToDos[0])
	}
}
func Test_Read_404TitleNotFound(t *testing.T) {

	initialToDos := []models.ToDo{
		{Title: "Clean room", Description: "Do laundry, hoover, clean sheets, change bin", Due: time.Date(2024, 9, 15, 12, 0, 0, 0, time.Local), Priority: 2, Status: "not started"},
		{Title: "Clean car", Description: "Remove mess, hoover, wash outside", Due: time.Date(2024, 9, 12, 12, 0, 0, 0, time.Local), Priority: 4, Status: "not started"}}
	server := Server{&sync.RWMutex{}, initialToDos, StubDecodeToDo(models.ToDo{}, nil)}
	response := httptest.NewRecorder()
	c := MockCustomContext("GET", models.ToDoRequest{}, []gin.Param{{Key: "title", Value: "slackOff"}}, response)

	server.Read(c)

	if response.Code != 404 {
		t.Error("Expected response to be 404 but was actually", response.Code)
	}
}
func Test_Delete_HappyPath(t *testing.T) {

	initialToDos := []models.ToDo{
		{Title: "Clean room", Description: "Do laundry, hoover, clean sheets, change bin", Due: time.Date(2024, 9, 15, 12, 0, 0, 0, time.Local), Priority: 2, Status: "not started"},
		{Title: "Clean car", Description: "Remove mess, hoover, wash outside", Due: time.Date(2024, 9, 12, 12, 0, 0, 0, time.Local), Priority: 4, Status: "not started"}}
	server := Server{&sync.RWMutex{}, initialToDos, StubDecodeToDo(models.ToDo{}, nil)}
	response := httptest.NewRecorder()
	c := MockCustomContext("POST", models.ToDoRequest{}, []gin.Param{{Key: "title", Value: "CleanCar"}}, response)

	server.Delete(c)

	if response.Code != 200 {
		t.Error("Expected response to be 200 but was actually", response.Code)
	}
	if l := len(server.StoredToDos); l != 1 {
		t.Error("Expected only 1 to-do to remain after deletion but was ", l)
	}
	if s := server.StoredToDos[0]; s != initialToDos[0] {
		t.Error("Expected remaining to-do", s, "to equal", initialToDos[0])
	}
}
func Test_Delete_404TitleNotFound(t *testing.T) {

	initialToDos := []models.ToDo{
		{Title: "Clean room", Description: "Do laundry, hoover, clean sheets, change bin", Due: time.Date(2024, 9, 15, 12, 0, 0, 0, time.Local), Priority: 2, Status: "not started"},
		{Title: "Clean car", Description: "Remove mess, hoover, wash outside", Due: time.Date(2024, 9, 12, 12, 0, 0, 0, time.Local), Priority: 4, Status: "not started"}}
	server := Server{&sync.RWMutex{}, initialToDos, StubDecodeToDo(models.ToDo{}, nil)}
	response := httptest.NewRecorder()
	c := MockCustomContext("POST", models.ToDoRequest{}, []gin.Param{{Key: "title", Value: "takeItEasy"}}, response)

	server.Delete(c)

	if response.Code != 404 {
		t.Error("Expected response to be 404 but was actually", response.Code)
	}
}
func Test_Update_HappyPath(t *testing.T) {

	initialToDos := []models.ToDo{
		{Title: "Clean room", Description: "Do laundry, hoover, clean sheets, change bin", Due: time.Date(2024, 9, 15, 12, 0, 0, 0, time.Local), Priority: 2, Status: "not started"},
		{Title: "Clean car", Description: "Remove mess, hoover, wash outside", Due: time.Date(2024, 9, 12, 12, 0, 0, 0, time.Local), Priority: 4, Status: "not started"}}
	expected := models.ToDo{Title: "Clean room", Description: "an updated description", Due: time.Date(2024, 9, 15, 12, 0, 0, 0, time.Local), Priority: 2, Status: "not started"}
	server := Server{&sync.RWMutex{}, initialToDos, StubDecodeToDo(expected, nil)}
	response := httptest.NewRecorder()
	c := MockCustomContext("POST", models.ToDoRequest{Title: "Clean room", Description: "an updated description", Due: "2024 9 15 12", Priority: 2, Status: "not started"}, []gin.Param{{Key: "title", Value: "Clean Room"}}, response)

	server.Update(c)

	if response.Code != 200 {
		t.Error("Expected response to be 200 but was actually", response.Code)
	}
	if l := len(server.StoredToDos); l != 2 {
		t.Error("Expected 2 to-dos to remain after update but was ", l)
	}
	if s := server.StoredToDos[0]; s != expected {
		t.Error("Expected remaining to-do", s, "to equal", expected)
	}
	if s := server.StoredToDos[1]; s != initialToDos[1] {
		t.Error("Expected remaining to-do", s, "to equal", initialToDos[1])
	}
}
func Test_Update_404TitleNotFound(t *testing.T) {

	initialToDos := []models.ToDo{
		{Title: "Clean room", Description: "Do laundry, hoover, clean sheets, change bin", Due: time.Date(2024, 9, 15, 12, 0, 0, 0, time.Local), Priority: 2, Status: "not started"},
		{Title: "Clean car", Description: "Remove mess, hoover, wash outside", Due: time.Date(2024, 9, 12, 12, 0, 0, 0, time.Local), Priority: 4, Status: "not started"}}
	expected := models.ToDo{Title: "Clean room", Description: "an updated description", Due: time.Date(2024, 9, 15, 12, 0, 0, 0, time.Local), Priority: 2, Status: "not started"}
	server := Server{&sync.RWMutex{}, initialToDos, StubDecodeToDo(expected, nil)}
	response := httptest.NewRecorder()
	c := MockCustomContext("POST", models.ToDoRequest{Title: "Clean room", Description: "an updated description", Due: "2024 9 15 12", Priority: 2, Status: "not started"}, []gin.Param{{Key: "title", Value: "SnoozeAlarm"}}, response)

	server.Update(c)

	if response.Code != 404 {
		t.Error("Expected response to be 404 but was actually", response.Code)
	}
}
func Test_Update_409TitleAlreadyExists(t *testing.T) {

	initialToDos := []models.ToDo{
		{Title: "Clean room", Description: "Do laundry, hoover, clean sheets, change bin", Due: time.Date(2024, 9, 15, 12, 0, 0, 0, time.Local), Priority: 2, Status: "not started"},
		{Title: "Clean car", Description: "Remove mess, hoover, wash outside", Due: time.Date(2024, 9, 12, 12, 0, 0, 0, time.Local), Priority: 4, Status: "not started"}}
	expected := models.ToDo{Title: "Clean car", Description: "an updated description", Due: time.Date(2024, 9, 15, 12, 0, 0, 0, time.Local), Priority: 2, Status: "not started"}
	server := Server{&sync.RWMutex{}, initialToDos, StubDecodeToDo(expected, nil)}
	response := httptest.NewRecorder()
	c := MockCustomContext("POST", models.ToDoRequest{Title: "Clean Car", Description: "an updated description", Due: "2024 9 15 12", Priority: 2, Status: "not started"}, []gin.Param{{Key: "title", Value: "cleanRoom"}}, response)

	server.Update(c)

	if response.Code != 409 {
		t.Error("Expected response to be 409 but was actually", response.Code)
	}
}
func Test_Update_400BadRequest(t *testing.T) {

	initialToDos := []models.ToDo{
		{Title: "Clean room", Description: "Do laundry, hoover, clean sheets, change bin", Due: time.Date(2024, 9, 15, 12, 0, 0, 0, time.Local), Priority: 2, Status: "not started"},
		{Title: "Clean car", Description: "Remove mess, hoover, wash outside", Due: time.Date(2024, 9, 12, 12, 0, 0, 0, time.Local), Priority: 4, Status: "not started"}}
	server := Server{&sync.RWMutex{}, initialToDos, StubDecodeToDo(models.ToDo{}, errors.New("an error message"))}
	response := httptest.NewRecorder()
	c := MockCustomContext("POST", models.ToDoRequest{}, []gin.Param{}, response)

	server.Update(c)

	if response.Code != 400 {
		t.Error("Expected response to be 400 but was actually", response.Code)
	}
	if m := strings.Split(response.Body.String(), `"`)[3]; m != "an error message" {
		t.Error(`Expected response message to be "an error message" but was actually`, m)
	}
}
func Test_Create_HappyPath(t *testing.T) {

	initialToDos := []models.ToDo{
		{Title: "Clean room", Description: "Do laundry, hoover, clean sheets, change bin", Due: time.Date(2024, 9, 15, 12, 0, 0, 0, time.Local), Priority: 2, Status: "not started"},
		{Title: "Clean car", Description: "Remove mess, hoover, wash outside", Due: time.Date(2024, 9, 12, 12, 0, 0, 0, time.Local), Priority: 4, Status: "not started"}}
	expected := models.ToDo{Title: "new title", Description: "new description", Due: time.Date(2024, 1, 1, 1, 0, 0, 0, time.Local), Priority: 1, Status: "not started"}
	server := Server{&sync.RWMutex{}, initialToDos, StubDecodeToDo(expected, nil)}
	response := httptest.NewRecorder()
	c := MockCustomContext("POST", models.ToDoRequest{Title: "new title", Description: "new description", Due: "2024 1 1 1", Priority: 1, Status: "not started"}, []gin.Param{}, response)

	server.Create(c)

	if response.Code != 200 {
		t.Error("Expected response to be 200 but was actually", response.Code)
	}
	if l := len(server.StoredToDos); l != 3 {
		t.Error("Expected 3 to-dos to exist after creation but was ", l)
	}
	if s := server.StoredToDos[2]; s != expected {
		t.Error("Expected new to-do", s, "to equal", expected)
	}
	if s := server.StoredToDos[1]; s != initialToDos[1] {
		t.Error("Expected existing to-do", s, "to equal", initialToDos[1])
	}
	if s := server.StoredToDos[0]; s != initialToDos[0] {
		t.Error("Expected existing to-do", s, "to equal", initialToDos[0])
	}
}
func Test_Create_409TitleAlreadyExists(t *testing.T) {

	initialToDos := []models.ToDo{
		{Title: "Clean room", Description: "Do laundry, hoover, clean sheets, change bin", Due: time.Date(2024, 9, 15, 12, 0, 0, 0, time.Local), Priority: 2, Status: "not started"},
		{Title: "Clean car", Description: "Remove mess, hoover, wash outside", Due: time.Date(2024, 9, 12, 12, 0, 0, 0, time.Local), Priority: 4, Status: "not started"}}
	expected := models.ToDo{Title: "clean room", Description: "new description", Due: time.Date(2024, 1, 1, 1, 0, 0, 0, time.Local), Priority: 1, Status: "not started"}
	server := Server{&sync.RWMutex{}, initialToDos, StubDecodeToDo(expected, nil)}
	response := httptest.NewRecorder()
	c := MockCustomContext("POST", models.ToDoRequest{Title: "clean room", Description: "new description", Due: "2024 1 1 1", Priority: 1, Status: "not started"}, []gin.Param{}, response)

	server.Create(c)

	if response.Code != 409 {
		t.Error("Expected response to be 409 but was actually", response.Code)
	}
}
func Test_Create_400BadRequest(t *testing.T) {

	initialToDos := []models.ToDo{
		{Title: "Clean room", Description: "Do laundry, hoover, clean sheets, change bin", Due: time.Date(2024, 9, 15, 12, 0, 0, 0, time.Local), Priority: 2, Status: "not started"},
		{Title: "Clean car", Description: "Remove mess, hoover, wash outside", Due: time.Date(2024, 9, 12, 12, 0, 0, 0, time.Local), Priority: 4, Status: "not started"}}
	server := Server{&sync.RWMutex{}, initialToDos, StubDecodeToDo(models.ToDo{}, errors.New("missing data"))}
	response := httptest.NewRecorder()
	c := MockCustomContext("POST", models.ToDoRequest{}, []gin.Param{}, response)

	server.Create(c)

	if response.Code != 400 {
		t.Error("Expected response to be 400 but was actually", response.Code)
	}
	if m := strings.Split(response.Body.String(), `"`)[3]; m != "missing data" {
		t.Error(`Expected response message to be "missing data" but was actually`, m)
	}
}

func MockCustomContext(methodType string, tdr models.ToDoRequest, p []gin.Param, w *httptest.ResponseRecorder) *gin.Context {
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

func StubDecodeToDo(td models.ToDo, err error) func(c *gin.Context) (models.ToDo, error) {
	return func(c *gin.Context) (models.ToDo, error) {
		return td, err
	}
}
