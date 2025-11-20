# Go Hexagonal Architecture Boilerplate

A production-ready boilerplate for building scalable Go applications using hexagonal architecture (ports and adapters) with support for both REST and gRPC APIs, database migrations, worker services, and resilience patterns.

## Features

- **Hexagonal Architecture**: Clean separation of concerns with domain, usecase, repository, and delivery layers
- **Dual API Support**: REST API using Fiber and gRPC API from a single codebase
- **Database Support**: MySQL and PostgreSQL with GORM
- **Caching**: Redis support
- **Message Queue**: RabbitMQ integration
- **Database Migrations**: Using golang-migrate
- **Worker Services**: Separate worker processes with graceful shutdown
- **Resilience Patterns**: Circuit breaker, retry, timeout, bulkhead, rate limiting, deduplication
- **Structured Logging**: Using Zap logger
- **Configuration Management**: Using Viper
- **Input Validation**: Using Go Playground Validator
- **API Documentation**: Swagger/OpenAPI documentation
- **Containerization**: Docker support
- **CI/CD**: GitHub Actions workflows
- **Pagination**: Offset and cursor-based pagination helpers
- **Centralized Responses**: Unified API response format with metadata

## Project Structure

```
├── api/                    # API definitions (protobuf files)
├── cmd/                    # Application entry points
│   ├── main.go            # Main application entry point
│   └── worker.go          # Worker service entry point
├── config/                 # Configuration and bootstrapping
├── database/               # Database migrations
├── docs/                   # Documentation files
├── internal/               # Application core (hexagonal architecture layers)
│   ├── delivery/          # Delivery mechanisms (REST, gRPC)
│   ├── domain/            # Business domain models and interfaces
│   ├── helper/            # Helper functions
│   ├── middleware/        # Middleware functions
│   ├── repository/        # Database implementations
│   ├── resilience/        # Resilience patterns
│   └── usecase/           # Business logic
├── pkg/                    # Shared packages
│   ├── gorm/              # GORM helpers (pagination)
│   ├── mysql/             # MySQL connection
│   ├── postgres/          # PostgreSQL connection
│   ├── rabbitmq/          # RabbitMQ connection
│   └── redis/             # Redis connection
└── test/                   # Test files
```

## Getting Started

### Prerequisites

- Go 1.23+
- Docker (for database and message queue services)
- Docker Compose

### Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd go-hexagonal-boilerplate
   ```

2. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```

3. Install dependencies:
   ```bash
   go mod tidy
   ```

### Running the Application

#### Development Mode

Using Air for hot reloading:
```bash
make run-dev
```

#### Production Mode

```bash
# Build the application
make build

# Run the application
./app-hexagonal
```

#### Running Worker Services

```bash
# Run worker service
./app-hexagonal worker
```

### Database Migrations

```bash
# Run migrations
make migrate-up

# Rollback migrations
make migrate-down
```

### API Documentation

- REST API documentation: `http://localhost:4001/swagger/index.html`
- gRPC API: Defined in `api/proto/v1/`

## Features in Detail

### REST and gRPC APIs

The boilerplate supports both REST and gRPC APIs from a single codebase:
- REST API served on the configured HTTP port
- gRPC API served on the configured gRPC port
- Both APIs share the same business logic and data models

### Resilience Patterns

The application includes built-in resilience patterns:
- **Circuit Breaker**: Prevents cascading failures
- **Retry**: Automatic retry with exponential backoff
- **Timeout**: Prevents hanging requests
- **Bulkhead**: Limits concurrent requests
- **Rate Limiting**: Controls request rate
- **Deduplication**: Prevents duplicate processing

### Pagination

The boilerplate includes pagination helpers for GORM:
- Offset-based pagination
- Cursor-based pagination (simplified implementation)

### Centralized Responses

All API responses follow a consistent format:
```json
{
  "error": false,
  "code": 200,
  "message": "OK",
  "data": {},
  "metadata": {
    "timestamp": 1234567890,
    "page": {
      "current_page": 1,
      "page_size": 10,
      "total_records": 100,
      "total_pages": 10,
      "has_next": true,
      "has_previous": false
    }
  }
}
```

## Makefile Commands

```bash
make run-dev        # Run in development mode with hot reload
make build          # Build the application
make test           # Run tests
make migrate-up     # Run database migrations
make migrate-down   # Rollback database migrations
make proto-gen      # Generate protobuf files
make clean          # Clean build artifacts
```

## Configuration

The application is configured using environment variables. See `.env.example` for all available options.

## Docker

The application includes a Dockerfile for containerization:
```bash
docker build -t app-hexagonal .
docker run -p 4001:4001 app-hexagonal
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a pull request

## License

This project is licensed under the MIT License - see the [LICENSE.txt](LICENSE.txt) file for details.