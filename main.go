package main

import (
	"io"
	"net/http"

	"encoding/json"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

var db, _ = gorm.Open("mysql", "root:test1234@/todolist?charset=utf8&parseTime=True&loc=Local")

type Tasks struct {
	Id          int `gorm:"primary_key"`
	Name        string
	Description string
	Completed   bool
}

func ListTasks(w http.ResponseWriter, r *http.Request) {
	log.Info("Get all Tasks")
	tasks := GetTasks()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func GetTasks() interface{} {
	var tasks []Tasks
	results := db.Find(&tasks).Value
	return results
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	description := r.FormValue("description")
	log.WithFields(log.Fields{"description": description}).Info("Add new TodoItem. Saving to database.")
	task := &Tasks{Description: description, Completed: false}
	db.Create(&task)
	result := db.Last(&task)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result.Value)
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	log.Info("API Health is OK")
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive": true}`)
}

func handleRoutes() {
	router := mux.NewRouter()
	router.HandleFunc("/healthz", Healthz).Methods("GET")
	router.HandleFunc("/tasks", ListTasks).Methods("GET")
	router.HandleFunc("/tasks", CreateTask).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}

func main() {
	defer db.Close()

	db.Debug().DropTableIfExists(&Tasks{})
	db.Debug().AutoMigrate(&Tasks{})
	log.Info("Starting Task Management API server...")
	handleRoutes()
}
