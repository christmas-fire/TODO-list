package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"todo_list/internal/database"
	"todo_list/internal/tasks"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

// CORS middleware
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	db = database.InitDB()
	defer db.Close()

	r := mux.NewRouter()

	// Применяем CORS middleware
	r.Use(enableCORS)

	// Обработчик для корневого маршрута
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/web/templates/index.html")
	}).Methods("GET")

	// Обработчики для статических файлов
	r.HandleFunc("/static/styles.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "/static/styles.css")
	}).Methods("GET")

	r.HandleFunc("/static/script.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		http.ServeFile(w, r, "/static/script.js")
	}).Methods("GET")

	// API маршруты
	r.HandleFunc("/tasks", GetAllTasksHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/tasks/{id:[0-9]+}", GetTaskByIDHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/tasks", AddTaskHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/tasks/{id:[0-9]+}", UpdateTaskHandler).Methods("PATCH", "OPTIONS")
	r.HandleFunc("/tasks/{id:[0-9]+}", DeleteTaskHandler).Methods("DELETE", "OPTIONS")

	fmt.Println("Сервер запущен на http://localhost:8080")
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
	err = tasks.UpdateTask(db, id, t.Status, t.Title, t.Description)
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
