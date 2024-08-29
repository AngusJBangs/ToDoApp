package handlers

import (
	"ToDoApp/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (s *Server) ReadAll(c *gin.Context) {
	s.L.RLock()
	defer s.L.RUnlock()
	c.IndentedJSON(http.StatusOK, s.StoredToDos)
}
func (s *Server) Read(c *gin.Context) {
	title := c.Param("title")
	s.L.RLock()
	defer s.L.RUnlock()
	for _, td := range s.StoredToDos {
		if strings.EqualFold(strings.ReplaceAll(td.Title, " ", ""), strings.ReplaceAll(title, " ", "")) {
			c.IndentedJSON(http.StatusOK, td)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "to-do not found"})
}
func (s *Server) Create(c *gin.Context) {
	newToDo, err := s.DecodeToDo(c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	s.L.Lock()
	defer s.L.Unlock()
	for _, td := range s.StoredToDos {
		if strings.EqualFold(strings.ReplaceAll(td.Title, " ", ""), strings.ReplaceAll(newToDo.Title, " ", "")) {
			c.IndentedJSON(http.StatusConflict, gin.H{"message": "To do with that title already exists"})
			return
		}
	}
	s.StoredToDos = append(s.StoredToDos, newToDo)
	c.IndentedJSON(http.StatusOK, gin.H{"message": "to do added"})
}
func (s *Server) Update(c *gin.Context) {
	title := c.Param("title")
	newToDo, err := s.DecodeToDo(c)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	titleSame := strings.EqualFold(strings.ReplaceAll(newToDo.Title, " ", ""), strings.ReplaceAll(title, " ", ""))
	foundAt := -1
	var alreadyExists bool
	s.L.Lock()
	defer s.L.Unlock()
	for i, td := range s.StoredToDos {
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
	s.StoredToDos[foundAt] = newToDo
	c.IndentedJSON(http.StatusOK, newToDo)
}
func (s *Server) Delete(c *gin.Context) {
	title := c.Param("title")
	s.L.Lock()
	defer s.L.Unlock()
	tds := append([]models.ToDo{}, s.StoredToDos...)
	for i, td := range s.StoredToDos {
		if strings.EqualFold(strings.ReplaceAll(td.Title, " ", ""), strings.ReplaceAll(title, " ", "")) {
			before := tds[:i]
			after := []models.ToDo{}
			if i < len(tds) {
				after = tds[i+1:]
			}
			s.StoredToDos = append(before, after...)
			c.IndentedJSON(http.StatusOK, gin.H{"message": "to-do deleted"})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "to-do not found"})
}
