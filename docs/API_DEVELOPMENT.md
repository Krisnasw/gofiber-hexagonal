# API Development Guidelines

This document outlines the process for developing new APIs in this hexagonal/onion architecture project.

## Creating New APIs

When creating new APIs, follow these steps:

### 1. Define the Proto File

Create a new proto file in the `api/proto/v1` directory:

```protobuf
syntax = "proto3";

package v1;

option go_package = "app-hexagonal/api/v1";

// ServiceName represents the service
service ServiceName {
  // MethodName describes what the method does
  rpc MethodName(MethodNameRequest) returns (MethodNameResponse) {}
}
```

### 2. Generate Proto Files

Run the following command to generate Go code from the proto file:

```bash
make proto-gen
```

### 3. Create Domain Models

Define your domain models in the `internal/domain` directory if they don't already exist.

### 4. Create Use Case Interface

Define the use case interface in the appropriate file in `internal/usecase`:

```go
type ServiceNameUseCaseInterface interface {
    MethodName(params) (result, error)
}
```

### 5. Implement Use Case

Implement the use case in the use case struct:

```go
func (uc *ServiceNameUseCase) MethodName(params) (result, error) {
    // Implementation
}
```

### 6. Create Application Service

Create an application service in `internal/application`:

```go
type ServiceNameService struct {
    serviceNameUseCase usecase.ServiceNameUseCaseInterface
}

func NewServiceNameService(serviceNameUseCase usecase.ServiceNameUseCaseInterface) *ServiceNameService {
    return &ServiceNameService{
        serviceNameUseCase: serviceNameUseCase,
    }
}

func (s *ServiceNameService) MethodName(params) (result, error) {
    return s.serviceNameUseCase.MethodName(params)
}
```

### 7. Create gRPC Service

Create a gRPC service implementation in `internal/delivery/grpc`:

```go
type ServiceNameServiceServer struct {
    v1.UnimplementedServiceNameServiceServer
    serviceNameService *application.ServiceNameService
    logger *zap.Logger
}

func NewServiceNameServiceServer(serviceNameService *application.ServiceNameService, logger *zap.Logger) *ServiceNameServiceServer {
    return &ServiceNameServiceServer{
        serviceNameService: serviceNameService,
        logger: logger,
    }
}

func (s *ServiceNameServiceServer) MethodName(ctx context.Context, req *v1.MethodNameRequest) (*v1.MethodNameResponse, error) {
    // Implementation
}
```

### 8. Register gRPC Service

Register the service in the gRPC server in `internal/delivery/grpc/server.go`:

```go
serviceNameServiceServer := NewServiceNameServiceServer(serviceNameService, s.logger)
v1.RegisterServiceNameServiceServer(s.server, serviceNameServiceServer)
```

### 9. Update Main Function

Update the main function in `cmd/main.go` to create and pass the new services:

```go
serviceNameUsecase := usecase.NewServiceNameUsecase(dependencies)
serviceNameService := application.NewServiceNameService(serviceNameUsecase)
// Pass to gRPC server
```

## Best Practices

1. Always define interfaces for use cases to maintain loose coupling
2. Use the application layer as a facade for complex operations
3. Keep domain models clean and free of framework-specific code
4. Follow the single responsibility principle for each layer
5. Use dependency injection to maintain testability
6. Document all public APIs and methods
7. Write unit tests for all use cases and application services