FROM golang:1.24-alpine

WORKDIR /app

COPY services/booking-service/go.mod services/booking-service/go.sum ./
RUN go mod download

COPY services/booking-service/ ./

COPY services/booking-service/.env /app/.env
COPY services/booking-service/config /app/config

RUN go build -o booking-service ./cmd/main.go

CMD ["./booking-service"]
