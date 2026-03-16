# Clean Architecture: Orders - REST, gRPC e GraphQL

Sistema de pedidos implementado com Clean Architecture em Go, expondo um único Use Case através de três interfaces: REST, gRPC e GraphQL.

## Tecnologias

- Go 1.22
- Clean Architecture
- MySQL 8.0
- REST (chi router)
- gRPC
- GraphQL (gqlgen)
- Docker & Docker Compose

## Como Executar

```bash
docker compose up
```

O comando acima irá:
1. Subir o banco de dados MySQL
2. Aguardar o banco estar pronto
3. Executar as migrações automaticamente
4. Iniciar a aplicação com todos os serviços

## Portas dos Serviços

- **REST API**: http://localhost:8000
- **gRPC**: localhost:50051
- **GraphQL**: http://localhost:8080 (Playground) | http://localhost:8080/query (Endpoint)

## Testando

Use o arquivo `api.http` na raiz do projeto (extensão REST Client do VS Code).

### REST

**Criar Order:**
```bash
curl -X POST http://localhost:8000/order \
  -H "Content-Type: application/json" \
  -d '{"id":"a","price":100.5,"tax":0.5}'
```

**Listar Orders:**
```bash
curl http://localhost:8000/order
```

### GraphQL

Acesse http://localhost:8080 no navegador.

**Mutation - Criar Order:**
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

**Query - Listar Orders:**
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

Instale o grpcurl:
```bash
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

**Criar Order:**
```bash
grpcurl -plaintext -d '{"id":"c","price":150.00,"tax":15.00}' \
  localhost:50051 pb.OrderService/CreateOrder
```

**Listar Orders:**
```bash
grpcurl -plaintext -d '{}' \
  localhost:50051 pb.OrderService/ListOrders
```

## Arquitetura

```
├── cmd/server/              # Ponto de entrada
├── internal/
│   ├── entity/              # Entidades e interfaces
│   ├── usecase/             # Casos de uso (CreateOrder, ListOrders)
│   └── infra/
│       ├── database/        # Repository e migrações
│       ├── web/             # Handler REST
│       ├── grpc/            # Service gRPC
│       └── graph/           # Resolver GraphQL
├── api/                     # Protobuf definitions
├── migrations/              # SQL migrations
└── docker-compose.yaml
```

## Parar os Serviços

```bash
docker compose down
```

Para remover volumes:
```bash
docker compose down -v
```
