package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
	Name         string   `json:"name"` //другие поля задачи
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта
// ...

// GetTasks возвращает список всех задач
func GetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Кодируем задачи в JSON
	err := json.NewEncoder(w).Encode(tasks)
	if err != nil {
		// Если произошла ошибка кодирования, возвращаем статус 500 и сообщение об ошибке
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetTask возвращает конкретную задачу по ID
func GetTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	task, exists := tasks[id]

	// Обработка случая, когда задача не найдена
	if !exists {
		http.Error(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	// Устанавливаем заголовок Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Кодируем задачу в JSON и обрабатываем возможные ошибки
	if err := json.NewEncoder(w).Encode(task); err != nil {
		// При возникновении ошибки кодирования, возвращаем статус 500
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Создание новой задачи
func CreateTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Генерация уникального ID
	newTask.ID = uuid.NewString()

	tasks[newTask.ID] = newTask
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

// UpdateTask обновляет существующую задачу по ID
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var updatedTask Task
	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Проверяем, существует ли задача
	if _, exists := tasks[id]; !exists {
		http.Error(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	// Обновление задачи
	updatedTask.ID = id // повышает читаемость и узнаваемость кода для будущих доработок и изменения Task,
	// например если он будет содержать дополнительные поля
	tasks[id] = updatedTask

	// Установка заголовка Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Кодируем обновлённую задачу в JSON и обрабатываем возможные ошибки
	if err := json.NewEncoder(w).Encode(updatedTask); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// DeleteTask удаляет задачу по ID
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if _, exists := tasks[id]; !exists {
		http.Error(w, "Задача не найдена", http.StatusNotFound)
		return
	}
	delete(tasks, id)
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	// ...
	// Регистрируем обработчики
	r.Get("/tasks", GetTasks)
	r.Get("/tasks/{id}", GetTask)
	r.Post("/tasks", CreateTask)
	r.Put("/tasks/{id}", UpdateTask)
	r.Delete("/tasks/{id}", DeleteTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
