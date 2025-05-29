FROM golang:1.24-alpine

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o api-gateway ./cmd/main.go

CMD ["./api-gateway"]
