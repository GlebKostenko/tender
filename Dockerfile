# Используем официальный образ Go версии 1.23
FROM golang:1.23 AS builder

# Создаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum файлы
COPY go-tender-app/go.mod go-tender-app/go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем все остальные файлы проекта
COPY go-tender-app ./

# Собираем проект
RUN go build -o main .

# Используем более легкий базовый образ для выполнения
FROM debian:bullseye-slim

# Создаем рабочую директорию
WORKDIR /app

# Копируем собранный бинарный файл из стадии сборки
COPY --from=builder /app/main .

# Открываем порт, на котором будет работать приложение
EXPOSE 8080

# Команда для запуска приложения
CMD ["./main"]
