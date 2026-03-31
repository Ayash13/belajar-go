# Stage 1: Build the Go application
FROM golang:1.22-alpine AS builder

# Atur working directory di dalam container
WORKDIR /app

# Inisialisasi modul Go secara independen dan install dependency (pq)
COPY main.go .
RUN go mod init dockerapp
RUN go get github.com/lib/pq

# Build aplikasi menjadi binary bernama 'main'
RUN go build -o main .

# Stage 2: Create a lightweight image
FROM alpine:latest

WORKDIR /app

# Salin binary dari stage builder ke image yang bersih
COPY --from=builder /app/main .

# Expose port aplikasi
EXPOSE 8080

# Jalankan aplikasi
CMD ["./main"]
