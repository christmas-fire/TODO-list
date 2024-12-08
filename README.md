# TODO-List Application

## РГЗ по ИТ (СибГУТИ, 3 семестр) - fullstack приложение TODO-list

### Обзор
Fullstack приложение TODO-list, разработанное с использованием современных технологий для эффективного управления задачами.

### Особенности
- Создание, редактирование и удаление задач
- Категоризация и приоритезация задач
- Удобный веб-интерфейс

### Технологии
- **Backend**: Go (Golang)
- **Frontend**: HTML, CSS, JS
- **База данных**: PostgreSQL
- **Архитектура**: golang-standards/project-layout

### Установка и запуск

#### Prerequisites
- Go 1.20+
- Git

#### Шаги установки
1. Клонируйте репозиторий
```bash
git clone https://github.com/christmas-fire/TODO-list.git
cd todo-list
```

2. Установите зависимости
```bash
go mod download
```

3. Запустите приложение
```bash
go run cmd/app/main.go
```

### Структура проекта
```
todo-list/
│
├── cmd/           # Точка входа приложения
├── internal/      # Внутренняя логика приложения
│   ├── database/  # Работа с базой данных
│   ├── rest/      # HTTP-обработчики
│   └── tasks/     # Логика работы с задачами
└── web/           # Веб-интерфейс
    └── templates/ # HTML-шаблоны
```

---
*Создано с ❤️ для учебного проекта*
