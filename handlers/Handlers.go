package handlers

import (
	"ToDoApp/models"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type CommandType int

const (
	ReadAllCommand = iota
	ReadCommand
	UpdateCommand
	CreateCommand
	DeleteCommand
)

type indentedJsonResponse struct {
	httpstatus int
	response   any
}

type Command struct {
	ty        CommandType
	c         *gin.Context
	replyChan chan indentedJsonResponse
}

func (s *Server) InitiateToDoHandlerManager() chan<- Command {
	cmds := make(chan Command)

	go func() {
		for cmd := range cmds {
			switch cmd.ty {
			case ReadAllCommand:
				cmd.replyChan <- s.readAll()
			case ReadCommand:
				cmd.replyChan <- s.read(cmd.c)
			case UpdateCommand:
				cmd.replyChan <- s.update(cmd.c)
			case DeleteCommand:
				cmd.replyChan <- s.delete(cmd.c)
			case CreateCommand:
				cmd.replyChan <- s.create(cmd.c)
			}
		}
	}()

	return cmds
}

func (s *Server) StartReadAll(c *gin.Context) {
	replyChan := make(chan indentedJsonResponse)
	s.Cmds <- Command{ty: ReadAllCommand, c: c, replyChan: replyChan}
	fmt.Println("waiting for reply")
	response := <-replyChan
	c.IndentedJSON(response.httpstatus, response.response)
}
func (s *Server) readAll() indentedJsonResponse {
	return indentedJsonResponse{http.StatusOK, s.StoredToDos}
}
func (s *Server) StartRead(c *gin.Context) {
	replyChan := make(chan indentedJsonResponse)
	s.Cmds <- Command{ty: ReadCommand, c: c, replyChan: replyChan}
	response := <-replyChan
	c.IndentedJSON(response.httpstatus, response.response)
}
func (s *Server) read(c *gin.Context) indentedJsonResponse {
	title := c.Param("title")
	for _, td := range s.StoredToDos {
		if strings.EqualFold(strings.ReplaceAll(td.Title, " ", ""), strings.ReplaceAll(title, " ", "")) {
			return indentedJsonResponse{http.StatusOK, td}
		}
	}
	return indentedJsonResponse{http.StatusNotFound, gin.H{"message": "to-do not found"}}
}
func (s *Server) StartCreate(c *gin.Context) {
	replyChan := make(chan indentedJsonResponse)
	s.Cmds <- Command{ty: CreateCommand, c: c, replyChan: replyChan}
	response := <-replyChan
	c.IndentedJSON(response.httpstatus, response.response)
}
func (s *Server) create(c *gin.Context) indentedJsonResponse {
	newToDo, err := s.DecodeToDo(c)
	if err != nil {
		return indentedJsonResponse{http.StatusBadRequest, gin.H{"message": err.Error()}}
	}
	for _, td := range s.StoredToDos {
		if strings.EqualFold(strings.ReplaceAll(td.Title, " ", ""), strings.ReplaceAll(newToDo.Title, " ", "")) {
			return indentedJsonResponse{http.StatusConflict, gin.H{"message": "To do with that title already exists"}}
		}
	}
	s.StoredToDos = append(s.StoredToDos, newToDo)
	return indentedJsonResponse{http.StatusOK, gin.H{"message": "to do added"}}
}
func (s *Server) StartUpdate(c *gin.Context) {
	replyChan := make(chan indentedJsonResponse)
	s.Cmds <- Command{ty: UpdateCommand, c: c, replyChan: replyChan}
	response := <-replyChan
	c.IndentedJSON(response.httpstatus, response.response)
}
func (s *Server) update(c *gin.Context) indentedJsonResponse {
	title := c.Param("title")
	newToDo, err := s.DecodeToDo(c)
	if err != nil {
		return indentedJsonResponse{http.StatusBadRequest, gin.H{"message": err.Error()}}
	}
	titleSame := strings.EqualFold(strings.ReplaceAll(newToDo.Title, " ", ""), strings.ReplaceAll(title, " ", ""))
	foundAt := -1
	var alreadyExists bool
	for i, td := range s.StoredToDos {
		if !titleSame && strings.EqualFold(strings.ReplaceAll(td.Title, " ", ""), strings.ReplaceAll(newToDo.Title, " ", "")) {
			alreadyExists = true
		}
		if strings.EqualFold(strings.ReplaceAll(td.Title, " ", ""), strings.ReplaceAll(title, " ", "")) {
			foundAt = i
		}
	}
	if foundAt < 0 {
		return indentedJsonResponse{http.StatusNotFound, gin.H{"message": "to-do not found"}}
	}
	if alreadyExists {
		return indentedJsonResponse{http.StatusConflict, newToDo}
	}
	s.StoredToDos[foundAt] = newToDo
	return indentedJsonResponse{http.StatusOK, newToDo}
}
func (s *Server) StartDelete(c *gin.Context) {
	replyChan := make(chan indentedJsonResponse)
	s.Cmds <- Command{ty: DeleteCommand, c: c, replyChan: replyChan}
	response := <-replyChan
	c.IndentedJSON(response.httpstatus, response.response)
}
func (s *Server) delete(c *gin.Context) indentedJsonResponse {
	title := c.Param("title")
	tds := append([]models.ToDo{}, s.StoredToDos...)
	for i, td := range s.StoredToDos {
		if strings.EqualFold(strings.ReplaceAll(td.Title, " ", ""), strings.ReplaceAll(title, " ", "")) {
			before := tds[:i]
			after := []models.ToDo{}
			if i < len(tds) {
				after = tds[i+1:]
			}
			s.StoredToDos = append(before, after...)
			return indentedJsonResponse{http.StatusOK, gin.H{"message": "to-do deleted"}}
		}
	}
	return indentedJsonResponse{http.StatusNotFound, gin.H{"message": "to-do not found"}}
}
