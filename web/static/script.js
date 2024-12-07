// Глобальные переменные
let currentEditingId = null;

// Глобальные переменные для хранения текущих задач и фильтров
let allTasks = [];
let currentFilters = {
    search: '',
    status: 'all',
    date: 'all',
    sort: 'createDesc'
};

// Инициализация при загрузке страницы
document.addEventListener('DOMContentLoaded', () => {
    initializeFilters();
    initTheme();
    loadTasks();
    
    // Обработчики событий
    document.getElementById('themeToggle').addEventListener('click', toggleTheme);
    document.getElementById('newTaskBtn').addEventListener('click', () => showModal('new'));
    document.getElementById('saveTaskBtn').addEventListener('click', saveTask);
});

// Управление темой
function initTheme() {
    const theme = localStorage.getItem('theme') || 'light';
    document.documentElement.setAttribute('data-theme', theme);
}

function toggleTheme() {
    const currentTheme = document.documentElement.getAttribute('data-theme');
    const newTheme = currentTheme === 'light' ? 'dark' : 'light';
    document.documentElement.setAttribute('data-theme', newTheme);
    localStorage.setItem('theme', newTheme);
}

// Управление модальным окном
function showModal(mode, taskId = null, taskData = null) {
    const modal = document.getElementById('taskModal');
    const modalTitle = document.getElementById('modalTitle');
    const titleInput = document.getElementById('taskTitle');
    const descriptionInput = document.getElementById('taskDescription');
    
    currentEditingId = taskId;
    
    if (mode === 'new') {
        modalTitle.textContent = 'Новая задача';
        titleInput.value = '';
        descriptionInput.value = '';
    } else {
        modalTitle.textContent = 'Редактировать задачу';
        titleInput.value = taskData.title || '';
        descriptionInput.value = taskData.description || '';
    }
    
    modal.style.display = 'block';
}

function closeModal() {
    const modal = document.getElementById('taskModal');
    modal.style.display = 'none';
    currentEditingId = null;
}

// Сохранение задачи (создание новой или обновление существующей)
async function saveTask() {
    const titleInput = document.getElementById('taskTitle');
    const descriptionInput = document.getElementById('taskDescription');
    const title = titleInput.value.trim();
    const description = descriptionInput.value.trim();

    if (!title) {
        alert('Введите название задачи');
        return;
    }

    try {
        if (currentEditingId) {
            // Обновление существующей задачи
            const response = await fetch(`http://localhost:8080/tasks/${currentEditingId}`, {
                method: 'PATCH',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    title: title,
                    description: description,
                    status: false // Сохраняем текущий статус при редактировании
                })
            });

            if (!response.ok) {
                throw new Error('Ошибка при обновлении задачи');
            }
        } else {
            // Создание новой задачи
            const response = await fetch('http://localhost:8080/tasks', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    title: title,
                    description: description
                })
            });

            if (!response.ok) {
                throw new Error('Ошибка при создании задачи');
            }
        }

        closeModal();
        await loadTasks();
    } catch (error) {
        console.error('Ошибка:', error);
        alert(error.message);
    }
}

// Загрузка задач с применением фильтров
async function loadTasks() {
    try {
        const response = await fetch('http://localhost:8080/tasks');
        if (!response.ok) {
            throw new Error('Не удалось загрузить задачи');
        }
        allTasks = await response.json();
        applyFilters();
    } catch (error) {
        console.error('Ошибка:', error);
        alert(error.message);
    }
}

// Применение всех фильтров к задачам
function applyFilters() {
    let filteredTasks = [...allTasks];

    // Фильтр по поиску
    if (currentFilters.search) {
        const searchLower = currentFilters.search.toLowerCase();
        filteredTasks = filteredTasks.filter(task => 
            task.title.toLowerCase().includes(searchLower) ||
            task.description.toLowerCase().includes(searchLower)
        );
    }

    // Фильтр по статусу
    if (currentFilters.status !== 'all') {
        filteredTasks = filteredTasks.filter(task => 
            currentFilters.status === 'completed' ? task.status : !task.status
        );
    }

    // Фильтр по дате
    if (currentFilters.date !== 'all') {
        const now = new Date();
        const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
        const weekAgo = new Date(today.getTime() - 7 * 24 * 60 * 60 * 1000);
        const monthAgo = new Date(today.getTime() - 30 * 24 * 60 * 60 * 1000);

        filteredTasks = filteredTasks.filter(task => {
            const taskDate = new Date(task.create_time);
            switch (currentFilters.date) {
                case 'today':
                    return taskDate >= today;
                case 'week':
                    return taskDate >= weekAgo;
                case 'month':
                    return taskDate >= monthAgo;
                default:
                    return true;
            }
        });
    }

    // Сортировка
    filteredTasks.sort((a, b) => {
        const dateA = new Date(a.create_time);
        const dateB = new Date(b.create_time);
        const completeDateA = a.complete_time ? new Date(a.complete_time) : null;
        const completeDateB = b.complete_time ? new Date(b.complete_time) : null;

        switch (currentFilters.sort) {
            case 'createAsc':
                return dateA - dateB;
            case 'createDesc':
                return dateB - dateA;
            case 'completeAsc':
                if (!completeDateA && !completeDateB) return 0;
                if (!completeDateA) return 1;
                if (!completeDateB) return -1;
                return completeDateA - completeDateB;
            case 'completeDesc':
                if (!completeDateA && !completeDateB) return 0;
                if (!completeDateA) return 1;
                if (!completeDateB) return -1;
                return completeDateB - completeDateA;
            default:
                return dateB - dateA;
        }
    });

    displayTasks(filteredTasks);
}

// Инициализация обработчиков событий для фильтров
function initializeFilters() {
    // Поиск по названию
    const searchInput = document.getElementById('searchInput');
    searchInput.addEventListener('input', (e) => {
        currentFilters.search = e.target.value;
        applyFilters();
    });

    // Фильтр по статусу
    const statusFilter = document.getElementById('statusFilter');
    statusFilter.addEventListener('change', (e) => {
        currentFilters.status = e.target.value;
        applyFilters();
    });

    // Фильтр по дате
    const dateFilter = document.getElementById('dateFilter');
    dateFilter.addEventListener('change', (e) => {
        currentFilters.date = e.target.value;
        applyFilters();
    });

    // Сортировка
    const sortFilter = document.getElementById('sortFilter');
    sortFilter.addEventListener('change', (e) => {
        currentFilters.sort = e.target.value;
        applyFilters();
    });
}

// Отображение задач
function displayTasks(tasks) {
    const tasksList = document.getElementById('tasksList');
    tasksList.innerHTML = '';

    tasks.forEach(task => {
        const createDate = formatDate(task.create_time);
        const completeDate = task.complete_time ? formatDate(task.complete_time) : null;
        
        const taskElement = document.createElement('div');
        taskElement.className = `task-item ${task.status ? 'completed' : ''}`;
        taskElement.innerHTML = `
            <div class="task-content">
                <div class="task-title">${task.title}</div>
                <div class="task-description">${task.description || ''}</div>
                <div class="task-dates">
                    <small>Создано: ${createDate}</small>
                    ${task.status ? `<br><small>Выполнено: ${completeDate}</small>` : ''}
                </div>
            </div>
            <div class="task-actions">
                <button onclick="toggleStatus(${task.id}, ${!task.status})" class="status-btn">
                    ${task.status ? 'Отменить' : 'Выполнить'}
                </button>
                <button onclick="showModal('edit', ${task.id}, {
                    title: '${task.title.replace(/'/g, "\\'")}',
                    description: '${(task.description || '').replace(/'/g, "\\'")}'
                })" class="edit-btn">
                    Изменить
                </button>
                <button onclick="deleteTask(${task.id})" class="delete-btn">
                    Удалить
                </button>
            </div>
        `;
        tasksList.appendChild(taskElement);
    });
}

// Изменение статуса задачи
async function toggleStatus(id, newStatus) {
    try {
        // Сначала получаем текущую задачу
        const response = await fetch(`http://localhost:8080/tasks/${id}`);
        if (!response.ok) {
            throw new Error('Не удалось получить задачу');
        }
        const task = await response.json();

        // Обновляем статус, сохраняя остальные поля
        const updateResponse = await fetch(`http://localhost:8080/tasks/${id}`, {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                title: task.title,
                description: task.description,
                status: newStatus
            })
        });

        if (!updateResponse.ok) {
            throw new Error('Не удалось обновить статус задачи');
        }

        await loadTasks();
    } catch (error) {
        console.error('Ошибка:', error);
        alert(error.message);
    }
}

// Удаление задачи
async function deleteTask(id) {
    if (!confirm('Вы уверены, что хотите удалить эту задачу?')) {
        return;
    }

    try {
        const response = await fetch(`http://localhost:8080/tasks/${id}`, {
            method: 'DELETE'
        });

        if (response.ok) {
            loadTasks();
        } else {
            throw new Error('Ошибка при удалении задачи');
        }
    } catch (error) {
        console.error('Ошибка:', error);
        alert('Не удалось удалить задачу');
    }
}

// Форматирование даты
function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleDateString('ru-RU', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
    });
}

// Закрытие модального окна при клике вне его
window.onclick = function(event) {
    const modal = document.getElementById('taskModal');
    if (event.target === modal) {
        closeModal();
    }
}
