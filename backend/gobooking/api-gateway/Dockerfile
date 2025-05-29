FROM golang:1.24-alpine

WORKDIR /app

COPY api-gateway/go.mod api-gateway/go.sum ./

RUN go mod download

COPY api-gateway/ ./

RUN go build -o api-gateway ./cmd/main.go

CMD ["./api-gateway"]
