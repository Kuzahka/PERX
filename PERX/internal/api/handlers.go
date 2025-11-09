package api

import (
	"PERX/internal/core"
	"PERX/internal/queue"
	"encoding/json"
	"net/http"
	"time"
)

type API struct {
	Storage    *queue.Storage
	Dispatcher *queue.Dispatcher
}

// PostTaskHandler - endpoint для постановки задачи в очередь
// POST /submit
func (a *API) PostTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	var params core.TaskParameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	if params.N < 1 {
		http.Error(w, "Параметр 'n' должен быть >= 1", http.StatusBadRequest)
		return
	}

	// 1. Добавляем в хранилище (получает статус "В очереди")
	taskID := a.Storage.AddTask(params)

	// 2. Отправляем задачу диспетчеру для обработки
	task, _ := a.Storage.GetTaskByID(taskID)
	a.Dispatcher.SubmitTask(task)

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"queue_id": taskID,
		"message":  "Задача успешно поставлена в очередь.",
	})
}

// GetTasksHandler - endpoint для получения списка задач
// GET /tasks
func (a *API) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	// 1. Получаем список из хранилища (включая проверку на TTL)
	tasks := a.Storage.GetTasks()

	// 2. Форматируем результат в нужный вид (поля результата)
	type TaskResult struct {
		ID               int             `json:"id"`
		Status           core.TaskStatus `json:"status"`
		N                int             `json:"n"`
		D                float64         `json:"d"`
		N1               float64         `json:"n1"`
		I                float64         `json:"i"`
		TTL              float64         `json:"ttl"`
		CurrentIteration int             `json:"current_iteration"`
		TimePlaced       time.Time       `json:"time_placed"`
		TimeStart        *time.Time      `json:"time_start,omitempty"`
		TimeEnd          *time.Time      `json:"time_end,omitempty"`
		CurrentValue     float64         `json:"current_value"`
	}

	results := make([]TaskResult, len(tasks))
	for i, t := range tasks {
		// Для N=1, итерация будет 0, а CurrentValue = N1
		currentIter := t.CurrentIteration
		if t.Status == core.StatusCompleted || t.Status == core.StatusExpired {
			currentIter = t.N
		}

		results[i] = TaskResult{
			ID:               t.ID,
			Status:           t.Status,
			N:                t.N,
			D:                t.D,
			N1:               t.N1,
			I:                t.I,
			TTL:              t.TTL,
			CurrentIteration: currentIter,
			TimePlaced:       t.TimePlaced,
			TimeStart:        t.TimeStart,
			TimeEnd:          t.TimeEnd,
			CurrentValue:     t.CurrentValue,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
