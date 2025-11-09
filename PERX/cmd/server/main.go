package main

import (
	"PERX/internal/api"
	"PERX/internal/queue"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func main() {
	// Получаем N из параметров командной строки
	if len(os.Args) < 2 {
		fmt.Println("Использование: server <количество_параллельных_задач_N>")
		os.Exit(1)
	}

	maxWorkers, err := strconv.Atoi(os.Args[1])
	if err != nil || maxWorkers < 1 {
		fmt.Println("Неверное значение для N. Должно быть целым числом > 0.")
		os.Exit(1)
	}

	fmt.Printf("Запуск сервера с максимальным количеством параллельных задач N = %d\n", maxWorkers)

	// Инициализация компонентов
	storage := queue.NewStorage()
	dispatcher := queue.NewDispatcher(storage, maxWorkers)

	// Запуск диспетчера и воркеров
	dispatcher.Start()

	// Инициализация API
	apiHandler := api.API{
		Storage:    storage,
		Dispatcher: dispatcher,
	}

	// Настройка роутов
	http.HandleFunc("/submit", apiHandler.PostTaskHandler) // 1. Поставить задачу в очередь
	http.HandleFunc("/tasks", apiHandler.GetTasksHandler)  // 2. Получить отсортированный список задач

	// Запуск HTTP-сервера
	port := "8080"
	fmt.Printf("HTTP-сервер запущен на :%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %v\n", err)
		os.Exit(1)
	}
}
