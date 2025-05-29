FROM golang:1.24-alpine

WORKDIR /app

COPY services/listing-service/go.mod services/listing-service/go.sum ./
RUN go mod download

COPY services/listing-service/ ./

RUN go build -o listing-service ./cmd/main.go

CMD ["./listing-service"]
