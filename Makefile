.PHONY: proto graphql generate run test clean docker-up docker-down

proto:
	@echo "Generating protobuf files..."
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/order.proto
	@echo "Protobuf files generated successfully!"

graphql:
	@echo "Generating GraphQL files..."
	@go run github.com/99designs/gqlgen generate
	@echo "GraphQL files generated successfully!"

generate: proto graphql
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Done!"

run:
	@go run cmd/server/main.go

test:
	@go test -v ./...

clean:
	@rm -rf internal/infra/grpc/pb/*.pb.go
	@rm -rf internal/infra/graph/generated.go
	@rm -rf internal/infra/graph/model
	@rm -rf internal/infra/graph/*.resolvers.go
	@echo "Cleaned generated files"

docker-up:
	@docker compose up --build

docker-down:
	@docker compose down

docker-down-v:
	@docker compose down -v

install-tools:
	@echo "Installing protoc plugins..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/99designs/gqlgen@latest
	@echo "Tools installed successfully!"

help:
	@echo "Available commands:"
	@echo "  make proto          - Generate protobuf files"
	@echo "  make graphql        - Generate GraphQL files"
	@echo "  make generate       - Generate proto, GraphQL and download dependencies"
	@echo "  make run            - Run application locally"
	@echo "  make test           - Run tests"
	@echo "  make clean          - Clean generated files"
	@echo "  make docker-up      - Start Docker containers"
	@echo "  make docker-down    - Stop Docker containers"
	@echo "  make docker-down-v  - Stop Docker containers and remove volumes"
	@echo "  make install-tools  - Install required tools"
