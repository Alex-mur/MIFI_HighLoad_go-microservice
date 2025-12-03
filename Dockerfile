FROM golang:1.24.3-alpine

WORKDIR /app

# Копируем go.mod
COPY go.mod .

# Скачиваем ВСЕ зависимости и создаем go.sum
RUN go mod download 2>/dev/null || (go mod init go-microservice && go mod tidy)

# Копируем весь исходный код
COPY . .

# Еще раз проверяем зависимости
RUN go mod tidy 2>/dev/null || true

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/server

EXPOSE 8080

CMD ["./app"]