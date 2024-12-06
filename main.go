package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"todo_list/database"
	"todo_list/tasks"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	db = database.InitDB()
	defer db.Close()

	http.HandleFunc("/tasks/get", GetAllTasksHandler)
	http.HandleFunc("/tasks/put", AddTaskHandler)

	fmt.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func GetAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := tasks.GetAllTasks(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t tasks.Task
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	id, err := tasks.AddTask(db, t.Title, t.Description)
	if err != nil {
		http.Error(w, "Ошибка при добавлении задачи", http.StatusInternalServerError)
		return
	}
	t.Id = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}
