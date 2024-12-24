# Этап 1: Сборка приложения
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Копируем модули и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальной код приложения
COPY . .

# Собираем приложение
RUN go build -o bin/crud_server cmd/main.go

# Этап 2: Минимальный образ для запуска приложения
FROM alpine:latest

WORKDIR /root/

# Копируем исполняемый файл из этапа сборки
COPY --from=builder /app/bin/crud_server .

CMD ["./crud_server"]
