package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func Test_Concurrency(t *testing.T) {
	initialToDos := MakeSampleToDos()

	errorMessages := make(chan string)

	requests, wg := GenerateRequests(initialToDos, errorMessages)
	time.Sleep(time.Second * 2)

	for _, r := range requests {
		go r()
	}
	go func() {
		defer close(errorMessages)
		wg.Wait()
	}()
	errorMessage := <-errorMessages
	if errorMessage != "" {
		t.Error("expected 200 response but recieved:", errorMessage)
	}
}

func GenerateRequests(toDos []ToDo, errs chan<- string) ([]func(), *sync.WaitGroup) {
	var requests []func()
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(5)
		requests = append(requests, ReadAllRequest(&wg, errs))
		requests = append(requests, ReadRequest(toDos[i].Title, &wg, errs))
		requests = append(requests, UpdateRequest(toDos[i].Title, "New description "+strconv.Itoa(i), &wg, errs))
		requests = append(requests, CreateRequest("New title "+strconv.Itoa(i+10), "New description "+strconv.Itoa(i+10), &wg, errs))
		requests = append(requests, DeleteRequest(toDos[i+5].Title, &wg, errs))
	}
	return requests, &wg
}

func ReadAllRequest(wg *sync.WaitGroup, errs chan<- string) func() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:8080/todos", nil)
	if err != nil {
		log.Fatal(err)
	}
	return func() {
		defer wg.Done()
		result, _ := client.Do(req)
		if result.Status != "200 OK" {
			errs <- result.Status
		}
	}
}

func ReadRequest(title string, wg *sync.WaitGroup, errs chan<- string) func() {
	client := &http.Client{}
	link := "http://localhost:8080/todo/" + strings.ReplaceAll(title, " ", "")
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		log.Fatal(err)
	}
	return func() {
		defer wg.Done()
		result, _ := client.Do(req)
		if result.Status != "200 OK" {
			errs <- result.Status
		}
	}
}
func UpdateRequest(title string, newDescription string, wg *sync.WaitGroup, errs chan<- string) func() {
	client := &http.Client{}
	data := strings.NewReader(`{"Title":"` + title + `","Description":"` + newDescription + `", "Due":"2024 08 29 17", "Priority":4,"Status":"not started"}`)
	link := "http://localhost:8080/update/" + strings.ReplaceAll(title, " ", "")
	req, err := http.NewRequest("POST", link, data)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Fatal(err)
	}
	return func() {
		defer wg.Done()
		result, _ := client.Do(req)
		if result.Status != "200 OK" {
			errs <- result.Status
		}
	}
}
func CreateRequest(newTitle string, newDescription string, wg *sync.WaitGroup, errs chan<- string) func() {
	client := &http.Client{}
	data := strings.NewReader(`{"Title":"` + newTitle + `","Description":"` + newDescription + `", "Due":"2024 08 29 17", "Priority":4,"Status":"not started"}`)
	link := "http://localhost:8080/create"
	req, err := http.NewRequest("POST", link, data)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Fatal(err)
	}
	return func() {
		defer wg.Done()
		result, _ := client.Do(req)
		if result.Status != "200 OK" {
			errs <- result.Status
		}
	}
}
func DeleteRequest(title string, wg *sync.WaitGroup, errs chan<- string) func() {
	client := &http.Client{}
	link := "http://localhost:8080/delete/" + strings.ReplaceAll(title, " ", "")
	req, err := http.NewRequest("POST", link, nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Fatal(err)
	}
	return func() {
		defer wg.Done()
		result, _ := client.Do(req)
		if result.Status != "200 OK" {
			errs <- result.Status
		}
	}

}
