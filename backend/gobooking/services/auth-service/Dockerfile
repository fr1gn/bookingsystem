FROM golang:1.24-alpine

WORKDIR /app

COPY services/auth-service/go.mod services/auth-service/go.sum ./
RUN go mod download

COPY services/auth-service/ ./

COPY services/auth-service/.env /app/.env
COPY services/auth-service/config /app/config

RUN go build -o auth-service ./cmd/main.go

CMD ["./auth-service"]
