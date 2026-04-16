FROM golang:alpine AS builder

RUN apk add --no-cache protobuf protobuf-dev
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

WORKDIR /app

COPY proto/ proto/
COPY pb/go.mod pb/go.mod

RUN mkdir -p pb && protoc --go_out=pb --go-grpc_out=pb \
  --go_opt=paths=source_relative \
  --go-grpc_opt=paths=source_relative \
  --proto_path=proto \
  proto/account.proto

COPY transaction-service/go.mod transaction-service/go.sum ./transaction-service/
WORKDIR /app/transaction-service
RUN go mod download

COPY transaction-service/ /app/transaction-service/

RUN go mod tidy && CGO_ENABLED=0 go build -o main .

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/transaction-service/main .
COPY --from=builder /app/transaction-service/.env .

EXPOSE 8082

CMD ["./main"]
