# Используем официальный образ Go как базовый
FROM golang:1.25-alpine as builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы с зависимостями
COPY go.mod go.sum ./

# Копируем исходники приложения в рабочую директорию
COPY . .

# Собираем приложение
RUN go build -o main ./cmd

# Начинаем новую стадию сборки на основе минимального образа
FROM alpine:latest

# Добавляем исполняемый файл, .env, migrations из первой стадии в корневую директорию контейнера
COPY --from=builder /app/main /main
COPY .env .env
COPY --from=builder /app/internal/infrastructure/migrations /app/internal/infrastructure/migrations

# Открываем порт 8080
EXPOSE 8080

# Запускаем приложение
CMD ["/main"]
