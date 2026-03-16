FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache protobuf-dev git

COPY go.mod go.sum ./
RUN go mod download

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

COPY . .

RUN protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    api/order.proto

RUN go run github.com/99designs/gqlgen generate

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8000 50051 8080

CMD ["./server"]
