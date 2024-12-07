package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"todo_list/database"
	"todo_list/tasks"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	db = database.InitDB()
	defer db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/tasks", GetAllTasksHandler).Methods("GET")
	r.HandleFunc("/tasks/{id:[0-9]+}", GetTaskByIDHandler).Methods("GET")
	r.HandleFunc("/tasks", AddTaskHandler).Methods("POST")
	r.HandleFunc("/tasks/{id:[0-9]+}", UpdateTaskHandler).Methods("PATCH")
	r.HandleFunc("/tasks/{id:[0-9]+}", DeleteTaskHandler).Methods("DELETE")

	fmt.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// Получить все задачи
func GetAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := tasks.GetAllTasks(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func GetTaskByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Некорректный ID задачи", http.StatusBadRequest)
		return
	}

	task, err := tasks.GetTaskByID(db, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// Добавить задачу
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

// Обновить задачу
func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Получение параметра {id} из URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Некорректный ID задачи", http.StatusBadRequest)
		return
	}

	var t tasks.Task
	err = json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Обновляем задачу в базе данных
	err = tasks.UpdateTask(db, id, t.Status, t.Description)
	if err != nil {
		http.Error(w, "Ошибка при обновлении задачи", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Удалить задачу
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Получение параметра {id} из URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Некорректный ID задачи", http.StatusBadRequest)
		return
	}

	// Удаляем задачу из базы данных
	err = tasks.DeleteTask(db, id)
	if err != nil {
		http.Error(w, "Ошибка при удалении задачи", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
