# Используем официальный образ Go
FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Копируем все исходные коды
COPY . .

RUN CGO_ENABLED=0 go build -o main ./cmd/web/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

COPY --from=builder /app/config/app.yaml ./config/app.yaml

RUN chmod +x main

CMD ["./main"]
