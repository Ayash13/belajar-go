#!/bin/bash
set -e

cd "$(dirname "$0")"

# Ensure grpc plugins, go, and protoc are in PATH
export PATH="/opt/homebrew/bin:/usr/local/bin:$PATH:$HOME/go/bin"

echo "==> Generating protobuf Go code..."

mkdir -p pb

protoc --go_out=pb --go-grpc_out=pb \
  --go_opt=paths=source_relative \
  --go-grpc_opt=paths=source_relative \
  --proto_path=proto \
  proto/account.proto

echo "==> Running go mod tidy for pb..."
cd pb && go mod tidy && cd ..

echo "==> Running go mod tidy for account-service..."
cd account-service && go mod tidy && cd ..

echo "==> Running go mod tidy for transaction-service..."
cd transaction-service && go mod tidy && cd ..

echo "==> Running go mod tidy for notification-service..."
cd notification-service && go mod tidy && cd ..

echo "==> Done!"
