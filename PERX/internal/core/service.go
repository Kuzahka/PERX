package core

import (
	"fmt"
	"time"
)

// ExecuteTaskStep выполняет один шаг арифметической прогрессии для задачи.
// Она блокирует выполнение до тех пор, пока не пройдёт необходимый интервал I
// с момента последней записи.
func ExecuteTaskStep(task *Task) {
	if task.Status != StatusInProgress {
		// Защита, чтобы не выполнять завершенные или не начатые задачи
		return
	}

	// 1. Проверка и ожидание интервала I
	requiredWait := time.Duration(task.I * float64(time.Second))
	elapsedSinceLastCalc := time.Since(task.LastCalculation)

	if elapsedSinceLastCalc < requiredWait {
		timeToWait := requiredWait - elapsedSinceLastCalc
		fmt.Printf("[Task %d] Ожидание %v для соблюдения интервала I\n", task.ID, timeToWait)
		time.Sleep(timeToWait)
	}

	// Обновляем время последней записи, чтобы следующее вычисление отсчитывалось от *этого* момента
	now := time.Now()
	task.LastCalculation = now

	// 2. Вычисление нового значения:
	// **Требование:** "Вычисление текущего значения должно высчитываться от предыдущего
	// значения по факту наступления времени (по интервалу), а не по формуле разницы
	// времени и количеству итераций."

	// Новое значение = Предыдущее значение + Дельта (D)
	newValue := task.CurrentValue + task.D

	// 3. Обновление состояния задачи
	task.CurrentValue = newValue
	task.CurrentIteration++
	task.ArithmeticResults = append(task.ArithmeticResults, newValue)

	fmt.Printf("[Task %d] Итерация %d/%d, Значение: %.2f\n",
		task.ID, task.CurrentIteration, task.N, newValue)

	// 4. Проверка на завершение
	if task.CurrentIteration >= task.N-1 {
		// Мы уже посчитали N1 (индекс 0) и N-1 элементов после него.
		// Если N=1, то CurrentIteration=0 и условие task.N-1==0. Прогрессия завершена.
		// Если N=5, то итерации 0, 1, 2, 3, 4. Мы завершаем после итерации 4.

		task.Status = StatusCompleted
		taskEndTime := time.Now()
		task.TimeEnd = &taskEndTime
		fmt.Printf("--- [Task %d] Завершена в %v ---\n", task.ID, taskEndTime.Format(time.Stamp))
	}
}
