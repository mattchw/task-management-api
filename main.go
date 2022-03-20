package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Tasks struct {
	ID              string `json:"id"`
	TaskName        string `json:"task_name"`
	TaskDescription string `json:"task_description"`
	Date            string `json:"date"`
}

var tasks []Tasks

func allTasks() {
	task1 := Tasks{
		ID:              "1",
		TaskName:        "Task 1",
		TaskDescription: "Task 1 Description",
		Date:            "2020-01-01",
	}
	tasks = append(tasks, task1)
	task2 := Tasks{
		ID:              "2",
		TaskName:        "Task 2",
		TaskDescription: "Task 2 Description",
		Date:            "2020-01-01",
	}
	tasks = append(tasks, task2)
	fmt.Println("Your tasks are: ", tasks)
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	taskId := mux.Vars(r)["id"]
	flag := false
	for _, task := range tasks {
		if task.ID == taskId {
			flag = true
			json.NewEncoder(w).Encode(task)
			break
		}
	}
	if !flag {
		json.NewEncoder(w).Encode(map[string]string{"error": "Task not found"})
	}
}

func handleRoutes() {
	router := mux.NewRouter()
	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/tasks/{id}", getTask).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	allTasks()
	fmt.Println("Task Management API...")
	handleRoutes()
}
