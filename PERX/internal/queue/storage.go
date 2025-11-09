package queue

import (
	"PERX/internal/core"
	"sort"
	"sync"
	"time"
)

type Storage struct {
	mu     sync.RWMutex
	tasks  map[int]*core.Task
	nextID int
}

func NewStorage() *Storage {
	return &Storage{
		tasks:  make(map[int]*core.Task),
		nextID: 1,
	}
}

// AddTask добавляет задачу и возвращает её ID
func (s *Storage) AddTask(params core.TaskParameters) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.nextID
	newTask := &core.Task{
		ID:                id,
		TaskParameters:    params,
		Status:            core.StatusInQueue,
		CurrentValue:      params.N1,
		TimePlaced:        time.Now(),
		LastCalculation:   time.Now(), // Инициализация для первой итерации
		ArithmeticResults: []float64{params.N1},
	}
	s.tasks[id] = newTask
	s.nextID++
	return id
}

// GetTasks возвращает отсортированный список всех задач
func (s *Storage) GetTasks() []*core.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var list []*core.Task
	for _, task := range s.tasks {
		// Проверяем TTL, но не удаляем, а помечаем
		if task.Status == core.StatusCompleted && time.Since(*task.TimeEnd).Seconds() > task.TTL {
			task.Status = core.StatusExpired
		}
		list = append(list, task)
	}

	// Сортировка по ID (номеру в очереди)
	sort.Slice(list, func(i, j int) bool {
		return list[i].ID < list[j].ID
	})

	return list
}

// GetTaskByID получает задачу по ID
func (s *Storage) GetTaskByID(id int) (*core.Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	task, ok := s.tasks[id]
	return task, ok
}

// UpdateTaskInPlace позволяет обработчику напрямую обновлять статус задачи
func (s *Storage) UpdateTaskInPlace(task *core.Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Так как передаётся указатель, достаточно просто заблокировать и разблокировать мьютекс
	// для обеспечения атомарности операции "обновления" (даже если это только изменение статуса/времени)
	// и чтобы изменения, сделанные обработчиком, были видны другим горутинам.
}
