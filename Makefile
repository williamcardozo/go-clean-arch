.PHONY: proto graphql generate run test clean docker-up docker-down docker-down-keep test-all

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
	@docker compose down -v

docker-down-keep:
	@docker compose down

install-tools:
	@echo "Installing protoc plugins..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/99designs/gqlgen@latest
	@echo "Tools installed successfully!"
	@echo ""
	@echo "Installing evans..."
	@brew tap ktr0731/evans
	@brew install evans || echo "evans already installed or brew not available"
	@echo "Done!"

help:
	@echo "Available commands:"
	@echo "  make proto          - Generate protobuf files"
	@echo "  make graphql        - Generate GraphQL files"
	@echo "  make generate       - Generate proto, GraphQL and download dependencies"
	@echo "  make run            - Run application locally"
	@echo "  make test           - Run tests"
	@echo "  make clean          - Clean generated files"
	@echo "  make docker-up      - Start Docker containers"
	@echo "  make docker-down    - Stop Docker containers and remove volumes (clean state)"
	@echo "  make docker-down-keep - Stop Docker containers keeping volumes (keep data)"
	@echo "  make install-tools  - Install required tools (protoc, gqlgen, evans)"

test-all:
	@echo "Testing all interfaces..."
	@echo "\n=== REST API - Create Order ==="
	@curl -s -X POST http://localhost:8000/order -H "Content-Type: application/json" -d '{"id":"rest-'$$(date +%s)'","price":100,"tax":10}' | jq . || true
	@echo "\n=== REST API - List Orders ==="
	@curl -s http://localhost:8000/order | jq 'length' | xargs -I {} echo "Total orders: {}"
	@echo "\n=== GraphQL - Create Order ==="
	@curl -s -X POST http://localhost:8080/query -H "Content-Type: application/json" -d '{"query":"mutation { createOrder(input: {id: \"gql-'$$(date +%s)'\", price: 200, tax: 20}) { id price tax finalPrice } }"}' | jq . || true
	@echo "\n=== GraphQL - List Orders ==="
	@curl -s -X POST http://localhost:8080/query -H "Content-Type: application/json" -d '{"query":"{ orders { id price tax finalPrice } }"}' | jq '.data.orders | length' | xargs -I {} echo "Total orders: {}"
	@echo "\n=== gRPC - Create Order ==="
	@echo '{"id":"grpc-'$$(date +%s)'","price":300,"tax":30}' | evans --host localhost --port 50051 -r cli call pb.OrderService.CreateOrder || true
	@echo "\n=== gRPC - List Orders ==="
	@echo '{}' | evans --host localhost --port 50051 -r cli call pb.OrderService.ListOrders | jq '.orders | length' | xargs -I {} echo "Total orders: {}"
