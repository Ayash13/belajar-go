FROM golang:alpine AS builder

WORKDIR /app

COPY notification-service/go.mod notification-service/go.sum ./notification-service/
WORKDIR /app/notification-service
RUN go mod download

COPY notification-service/ /app/notification-service/

RUN go mod tidy && CGO_ENABLED=0 go build -o main .

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/notification-service/main .
COPY --from=builder /app/notification-service/.env .

EXPOSE 8083

CMD ["./main"]
