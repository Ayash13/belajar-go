# Stage 1: Build the Go application
FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod edit -go=1.24
RUN go mod download

COPY . .

WORKDIR /app/project_bank_api
RUN go build -o main .

# Stage 2: Create a lightweight image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/project_bank_api/main .
COPY --from=builder /app/project_bank_api/.env .

EXPOSE 8081

CMD ["./main"]
