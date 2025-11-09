package core

import "time"

// TaskParameters - Параметры задачи при постановке в очередь
type TaskParameters struct {
	N   int     `json:"n"`   // количество элементов (целочисленное)
	D   float64 `json:"d"`   // дельта между элементами (вещественное)
	N1  float64 `json:"n1"`  // Стартовое значение (вещественное)
	I   float64 `json:"i"`   // Интервал между итерациями в секундах (вещественное)
	TTL float64 `json:"ttl"` // TTL - время хранения результата в секундах (вещественное)
}

// TaskStatus - Статус выполнения задачи
type TaskStatus string

const (
	StatusInQueue    TaskStatus = "В очереди"
	StatusInProgress TaskStatus = "В процессе"
	StatusCompleted  TaskStatus = "Завершена"
	StatusExpired    TaskStatus = "Истек TTL" // Для удобства фильтрации после завершения TTL
)

// Task - Полное состояние задачи
type Task struct {
	ID                int        // Номер в очереди (уникальный ID)
	TaskParameters               // Встраиваем параметры
	Status            TaskStatus // Текущий статус
	CurrentIteration  int        // Текущая итерация (от 0 до N-1)
	CurrentValue      float64    // Текущее значение
	TimePlaced        time.Time  // Время постановки задачи
	TimeStart         *time.Time // Время старта задачи
	TimeEnd           *time.Time // Время окончания задачи
	LastCalculation   time.Time  // Время последнего вычисления (для соблюдения интервала I)
	ArithmeticResults []float64  // Хранение всех вычисленных элементов
}
