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

For detailed documentation, see the main [README.md](../README.md) file.