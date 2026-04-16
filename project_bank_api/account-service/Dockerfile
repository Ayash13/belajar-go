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

COPY account-service/go.mod account-service/go.sum ./account-service/
WORKDIR /app/account-service
RUN go mod download

COPY account-service/ /app/account-service/

RUN go mod tidy && CGO_ENABLED=0 go build -o main .

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/account-service/main .
COPY --from=builder /app/account-service/.env .

EXPOSE 8081 50051

CMD ["./main"]
