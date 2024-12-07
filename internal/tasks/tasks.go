package tasks

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Task struct {
	Id           int        `json:"id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Status       bool       `json:"status"`
	CreateTime   time.Time  `json:"create_time"`
	CompleteTime *time.Time `json:"complete_time"`
}

func AddTask(db *sql.DB, title, description string) (int, error) {
	query := `
		INSERT INTO tasks (title, description)
		VALUES ($1, $2)
		RETURNING id;
	`

	var id int

	err := db.QueryRow(query, title, description).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func GetAllTasks(db *sql.DB) ([]Task, error) {
	// SQL-запрос для получения всех задач
	query := "SELECT id, title, description, status, create_time, complete_time FROM tasks"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Срез для хранения задач
	var tasks []Task

	// Пробегаем по результатам запроса
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.Id, &t.Title, &t.Description, &t.Status, &t.CreateTime, &t.CompleteTime)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	// Проверяем, если в запросе возникла ошибка
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func GetTaskByID(db *sql.DB, id int) (Task, error) {
	// SQL-запрос для получения задачи по ID
	query := "SELECT id, title, description, status, create_time, complete_time FROM tasks WHERE id = $1"
	row := db.QueryRow(query, id)

	var t Task
	err := row.Scan(&t.Id, &t.Title, &t.Description, &t.Status, &t.CreateTime, &t.CompleteTime)
	if err != nil {
		if err == sql.ErrNoRows {
			// Если задачи с таким ID нет, возвращаем ошибку
			return t, errors.New("задача не найдена")
		}
		return t, err
	}

	return t, nil
}

func UpdateTask(db *sql.DB, id int, status bool, title, description string) error {
	var completeTime interface{} = nil
	if status {
		now := time.Now()
		completeTime = now
	}

	// SQL-запрос для обновления задачи
	query := `
		UPDATE tasks
		SET status = $1, title = $2, description = $3, complete_time = $4
		WHERE id = $5;
	`

	// Выполняем запрос
	_, err := db.Exec(query, status, title, description, completeTime, id)
	if err != nil {
		return fmt.Errorf("не удалось обновить задачу: %v", err)
	}

	return nil
}

func DeleteTask(db *sql.DB, id int) error {
	// SQL-запрос для удаления задачи по ID
	query := `
		DELETE FROM tasks
		WHERE id = $1;
	`

	// Выполняем запрос
	_, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("не удалось удалить задачу: %v", err)
	}

	return nil
}

func GetTasksByStatus(db *sql.DB, status bool) ([]Task, error) {
	// SQL-запрос для фильтрации по статусу
	query := `
		SELECT id, title, description, status, create_time, complete_time
		FROM tasks
		WHERE status = $1;
	`
	rows, err := db.Query(query, status)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить задачи по статусу: %v", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.Id, &t.Title, &t.Description, &t.Status, &t.CreateTime, &t.CompleteTime)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func GetTasksByCreateDate(db *sql.DB, date string) ([]Task, error) {
	// SQL-запрос для фильтрации по дате создания
	query := `
		SELECT id, title, description, status, create_time, complete_time
		FROM tasks
		WHERE DATE(create_time) = $1;
	`
	rows, err := db.Query(query, date)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить задачи по дате создания: %v", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.Id, &t.Title, &t.Description, &t.Status, &t.CreateTime, &t.CompleteTime)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func GetTasksByKeyword(db *sql.DB, keyword string) ([]Task, error) {
	// SQL-запрос для поиска по ключевым словам
	query := `
		SELECT id, title, description, status, create_time, complete_time
		FROM tasks
		WHERE title ILIKE $1 OR description ILIKE $1;
	`
	rows, err := db.Query(query, "%"+keyword+"%")
	if err != nil {
		return nil, fmt.Errorf("не удалось найти задачи по ключевому слову: %v", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.Id, &t.Title, &t.Description, &t.Status, &t.CreateTime, &t.CompleteTime)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}
