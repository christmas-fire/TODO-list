package main

import (
	"fmt"
	"net/http"
	"todo_list/internal/database"
	"todo_list/internal/rest"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// Запуск базы данных
	db := database.InitDB()
	defer db.Close()

	// Инициализация маршрутизатора
	r := mux.NewRouter()

	// Настройка CORS
	r.Use(rest.EnableCORS)

	// Обработка статических ресурсов
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/templates/index.html")
	}).Methods("GET")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static/"))))

	// Обработка HTTP запросов
	r.HandleFunc("/tasks", rest.GetAllTasksHandler(db)).Methods("GET", "OPTIONS")
	r.HandleFunc("/tasks/{id:[0-9]+}", rest.GetTaskByIDHandler(db)).Methods("GET", "OPTIONS")
	r.HandleFunc("/tasks", rest.AddTaskHandler(db)).Methods("POST", "OPTIONS")
	r.HandleFunc("/tasks/{id:[0-9]+}", rest.UpdateTaskHandler(db)).Methods("PATCH", "OPTIONS")
	r.HandleFunc("/tasks/{id:[0-9]+}", rest.DeleteTaskHandler(db)).Methods("DELETE", "OPTIONS")

	// Запуск сервера
	fmt.Println("Сервер запущен на http://localhost:8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
		return
	}
}
