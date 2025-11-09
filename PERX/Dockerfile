# Этап 1: Сборка
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение, статический бинарник для минимального конечного образа
RUN CGO_ENABLED=0 go build -o /server ./cmd/server/main.go

# Этап 2: Финальный образ
FROM alpine:latest

# Установка часового пояса (опционально, но полезно для корректного времени)
RUN apk --no-cache add tzdata

WORKDIR /root/

# Копируем собранный бинарник из builder'а
COPY --from=builder /server .

# Копируем часовой пояс
COPY --from=builder /usr/share/zoneinfo/Etc/UTC /usr/share/zoneinfo/Etc/UTC
ENV TZ=Etc/UTC

# Открываем порт
EXPOSE 8080

# Запуск приложения с параметром N (количество параллельных задач)
# ВАЖНО: <N> необходимо будет передать при запуске контейнера
# Пример: docker run -p 8080:8080 my-arithmetic-service 5
ENTRYPOINT ["./server"]
CMD ["1"] # Значение по умолчанию: 1 воркер