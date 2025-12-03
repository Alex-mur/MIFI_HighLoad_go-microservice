# Build stage
FROM golang:1.24.2-alpine AS builder

WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod ./
RUN go mod tidy
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарный файл
COPY --from=builder /app/main .

# Экспонируем порт
EXPOSE 8080

# Переменные окружения
ENV LOG_LEVEL=info
ENV GIN_MODE=release

# Запуск приложения
CMD ["./main"]