package queue

import (
	"PERX/internal/core"
	"fmt"
	"time"
)

// Dispatcher управляет пулом исполнителей (воркеров)
type Dispatcher struct {
	storage    *Storage
	taskQueue  chan *core.Task      // Канал для задач "в очереди"
	workerPool chan chan *core.Task // Канал для каналов воркеров (для выбора следующей задачи)
	maxWorkers int
	quit       chan struct{}
}

// NewDispatcher создает новый диспетчер с N воркерами
func NewDispatcher(storage *Storage, maxWorkers int) *Dispatcher {
	return &Dispatcher{
		storage:    storage,
		taskQueue:  make(chan *core.Task),
		workerPool: make(chan chan *core.Task, maxWorkers),
		maxWorkers: maxWorkers,
		quit:       make(chan struct{}),
	}
}

// Start запускает воркеры и начинает диспетчеризацию
func (d *Dispatcher) Start() {
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(i+1, d.workerPool, d.storage)
		worker.Start()
	}

	go d.dispatch()
}

// dispatch - главная горутина, которая выбирает следующую задачу
func (d *Dispatcher) dispatch() {
	for {
		select {
		case task := <-d.taskQueue:
			// Получена новая задача:
			go func(task *core.Task) {
				// Ждем, пока освободится канал воркера
				workerChannel := <-d.workerPool
				// Отправляем задачу свободному воркеру
				workerChannel <- task
			}(task)

		case <-d.quit:
			// Сигнал остановки
			return
		}
	}
}

// SubmitTask отправляет задачу в очередь диспетчеру
func (d *Dispatcher) SubmitTask(task *core.Task) {
	d.taskQueue <- task
}

// Worker - структура воркера, который выполняет задачи
type Worker struct {
	ID          int
	workerPool  chan chan *core.Task
	taskChannel chan *core.Task // Канал, на котором этот воркер ждет задач
	storage     *Storage
	quit        chan struct{}
}

func NewWorker(id int, workerPool chan chan *core.Task, storage *Storage) *Worker {
	return &Worker{
		ID:          id,
		workerPool:  workerPool,
		taskChannel: make(chan *core.Task),
		storage:     storage,
		quit:        make(chan struct{}),
	}
}

// Start запускает горутину воркера
func (w *Worker) Start() {
	go func() {
		for {
			// 1. Регистрация: Воркер сообщает диспетчеру, что он свободен.
			w.workerPool <- w.taskChannel

			select {
			case task := <-w.taskChannel:
				// 2. Получена задача, начинаем работу

				// Устанавливаем статус "В процессе" и время старта
				task.Status = core.StatusInProgress
				startTime := time.Now()
				task.TimeStart = &startTime
				w.storage.UpdateTaskInPlace(task) // Уведомить хранилище об изменении

				fmt.Printf("[Worker %d] Начал обработку задачи %d\n", w.ID, task.ID)

				// Выполняем итерации пока задача не завершится
				for task.Status != core.StatusCompleted {
					core.ExecuteTaskStep(task)
					w.storage.UpdateTaskInPlace(task) // Уведомить хранилище об изменении
				}

				fmt.Printf("[Worker %d] Завершил задачу %d\n", w.ID, task.ID)

			case <-w.quit:
				// 3. Сигнал остановки
				return
			}
		}
	}()
}
