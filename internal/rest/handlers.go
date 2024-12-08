package rest

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"todo_list/internal/tasks"

	"github.com/gorilla/mux"
)

func EnableCORS(next http.Handler) http.Handler {
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

func GetAllTasksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tasks, err := tasks.GetAllTasks(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)
	}
}

func GetTaskByIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}

func AddTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}

func UpdateTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		err = tasks.UpdateTask(db, id, t.Status, t.Title, t.Description)
		if err != nil {
			http.Error(w, "Ошибка при обновлении задачи", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func DeleteTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Некорректный ID задачи", http.StatusBadRequest)
			return
		}

		err = tasks.DeleteTask(db, id)
		if err != nil {
			http.Error(w, "Ошибка при удалении задачи", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
