FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o todo-api ./cmd/todo/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/todo-api .

CMD ["./todo-api"]