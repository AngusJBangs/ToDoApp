package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (s *Server) readAll(c *gin.Context) {
	s.l.RLock()
	defer s.l.RUnlock()
	c.IndentedJSON(http.StatusOK, s.storedToDos)
}
func (s *Server) read(c *gin.Context) {
	title := c.Param("title")
	s.l.RLock()
	defer s.l.RUnlock()
	for _, td := range s.storedToDos {
		if strings.EqualFold(strings.ReplaceAll(td.Title, " ", ""), strings.ReplaceAll(title, " ", "")) {
			c.IndentedJSON(http.StatusOK, td)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "to-do not found"})
}
func (s *Server) create(c *gin.Context) {
	newToDo, err := s.DecodeToDo(c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	s.l.Lock()
	defer s.l.Unlock()
	for _, td := range s.storedToDos {
		if strings.EqualFold(strings.ReplaceAll(td.Title, " ", ""), strings.ReplaceAll(newToDo.Title, " ", "")) {
			c.IndentedJSON(http.StatusConflict, gin.H{"message": "To do with that title already exists"})
			return
		}
	}
	s.storedToDos = append(s.storedToDos, newToDo)
	c.IndentedJSON(http.StatusOK, gin.H{"message": "to do added"})
}
func (s *Server) update(c *gin.Context) {
	title := c.Param("title")
	newToDo, err := s.DecodeToDo(c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	titleSame := strings.EqualFold(strings.ReplaceAll(newToDo.Title, " ", ""), strings.ReplaceAll(title, " ", ""))
	foundAt := -1
	var alreadyExists bool
	s.l.Lock()
	defer s.l.Unlock()
	for i, td := range s.storedToDos {
		if !titleSame && strings.EqualFold(strings.ReplaceAll(td.Title, " ", ""), strings.ReplaceAll(newToDo.Title, " ", "")) {
			alreadyExists = true
		}
		if strings.EqualFold(strings.ReplaceAll(td.Title, " ", ""), strings.ReplaceAll(title, " ", "")) {
			foundAt = i
		}
	}
	if foundAt < 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "to-do not found"})
		return
	}
	if alreadyExists {
		c.IndentedJSON(http.StatusConflict, newToDo)
		return
	}
	s.storedToDos[foundAt] = newToDo
	c.IndentedJSON(http.StatusOK, newToDo)
}
func (s *Server) delete(c *gin.Context) {
	title := c.Param("title")
	s.l.Lock()
	defer s.l.Unlock()
	tds := append([]ToDo{}, s.storedToDos...)
	for i, td := range s.storedToDos {
		if strings.EqualFold(strings.ReplaceAll(td.Title, " ", ""), strings.ReplaceAll(title, " ", "")) {
			before := tds[:i]
			after := []ToDo{}
			if i < len(tds) {
				after = tds[i+1:]
			}
			s.storedToDos = append(before, after...)
			c.IndentedJSON(http.StatusOK, gin.H{"message": "to-do deleted"})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "to-do not found"})
}
