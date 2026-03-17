# Clean Architecture: Orders - REST, gRPC e GraphQL

Sistema de pedidos com Clean Architecture em Go. Mesmos use cases expostos via REST, gRPC e GraphQL.

## Stack

- Go 1.25.0
- Clean Architecture
- MySQL 8.0
- REST (chi v5.0.11)
- gRPC (v1.79.2)
- GraphQL (gqlgen v0.17.88)
- Docker & Docker Compose

## Rodar o projeto

```bash
docker compose up
```

Sobe o MySQL, roda as migrations e inicia a aplicação.

## Portas

- **REST API**: http://localhost:8000
- **gRPC**: localhost:50051
- **GraphQL**: http://localhost:8080 (Playground) | http://localhost:8080/query (Endpoint)

## Testes

Arquivo `api.http` na raiz (REST Client do VS Code) ou:

```bash
make test-all
```

### REST

Criar Order:
```bash
curl -X POST http://localhost:8000/order \
  -H "Content-Type: application/json" \
  -d '{"id":"a","price":100.5,"tax":0.5}'
```

Listar Orders:
```bash
curl http://localhost:8000/order
```

### GraphQL

Acesse http://localhost:8080

Mutation - Criar Order:
```graphql
mutation {
  createOrder(input: {id: "b", price: 200.00, tax: 20.00}) {
    id
    price
    tax
    finalPrice
  }
}
```

Query - Listar Orders:
```graphql
{
  orders {
    id
    price
    tax
    finalPrice
  }
}
```

### gRPC

Instalar evans:
```bash
brew tap ktr0731/evans
brew install evans
```

Modo interativo:
```bash
evans -r repl -p 50051
```

Dentro do evans:
- `package pb`
- `service OrderService`
- `call CreateOrder`
- `call ListOrders`
- `show service`
- `exit`

Modo CLI - Criar Order:
```bash
echo '{"id":"c","price":300,"tax":30}' | evans --host localhost --port 50051 -r cli call pb.OrderService.CreateOrder
```

Modo CLI - Listar Orders:
```bash
echo '{}' | evans --host localhost --port 50051 -r cli call pb.OrderService.ListOrders
```

## Arquitetura

```
├── cmd/server/              # Entrypoint
├── internal/
│   ├── entity/              # Entities
│   ├── usecase/             # Use cases
│   └── infra/
│       ├── database/        # Repository
│       ├── web/             # REST handler
│       ├── grpc/            # gRPC service
│       └── graph/           # GraphQL resolver
├── api/                     # Protobuf
├── migrations/              # SQL migrations
└── docker-compose.yaml
```

## Parar

```bash
make docker-down
```

Remove containers e volumes (banco zerado).

Manter dados:
```bash
make docker-down-keep
```

## Comandos Make

```bash
make help              # Lista comandos
make install-tools     # Instala ferramentas
make generate          # Gera proto + GraphQL
make docker-up         # Sobe containers
make docker-down       # Para e limpa volumes
make docker-down-keep  # Para mantendo volumes
make test-all          # Testa tudo
```
