# Stage 1: Build the Go application
FROM golang:alpine AS builder

WORKDIR /app

# Salin go.mod dan go.sum dari parent (root project)
COPY go.mod go.sum ./
RUN go mod edit -go=1.24
RUN go mod download

# Salin seluruh source code root
COPY . .

# Pindah ke sub direktory target dan build
WORKDIR /app/challenge_3
RUN go build -o main .

# Stage 2: Create a lightweight image
FROM alpine:latest

WORKDIR /app

# Salin hasil build binary dari stage builder
COPY --from=builder /app/challenge_3/main .

EXPOSE 8080

CMD ["./main"]
