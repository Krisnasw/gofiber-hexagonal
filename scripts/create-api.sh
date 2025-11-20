#!/bin/bash

# Script to create a new API with proto file
# Usage: ./scripts/create-api.sh <service_name>

if [ $# -eq 0 ]; then
    echo "Usage: $0 <service_name>"
    echo "Example: $0 product"
    exit 1
fi

SERVICE_NAME=$1
PROTO_FILE="api/proto/v1/${SERVICE_NAME}.proto"

# Create proto file
cat > $PROTO_FILE << EOF
syntax = "proto3";

package v1;

option go_package = "app-hexagonal/api/v1";

// ${SERVICE_NAME^}Service represents the ${SERVICE_NAME} service
service ${SERVICE_NAME^}Service {
  // Get${SERVICE_NAME^} retrieves a ${SERVICE_NAME} by ID
  rpc Get${SERVICE_NAME^}(Get${SERVICE_NAME^}Request) returns (Get${SERVICE_NAME^}Response) {}
  
  // Create${SERVICE_NAME^} creates a new ${SERVICE_NAME}
  rpc Create${SERVICE_NAME^}(Create${SERVICE_NAME^}Request) returns (Create${SERVICE_NAME^}Response) {}
}
EOF

echo "Created proto file: $PROTO_FILE"

# Generate Go code from proto
make proto-gen

echo "Generated Go code from proto files"

# Create usecase interface and implementation
USECASE_FILE="internal/usecase/${SERVICE_NAME}.go"
cat > $USECASE_FILE << EOF
package usecase

import "app-hexagonal/internal/domain"

// ${SERVICE_NAME^}UsecaseInterface defines the interface for ${SERVICE_NAME} use cases
type ${SERVICE_NAME^}UsecaseInterface interface {
	Get${SERVICE_NAME^}ByID(id string) (*domain.${SERVICE_NAME^}, error)
	Create${SERVICE_NAME^}(${SERVICE_NAME} *domain.${SERVICE_NAME^}) error
}

// ${SERVICE_NAME^}Usecase handles ${SERVICE_NAME} business logic
type ${SERVICE_NAME^}Usecase struct {
	repo domain.${SERVICE_NAME^}Repository
}

// New${SERVICE_NAME^}Usecase creates a new ${SERVICE_NAME} usecase
func New${SERVICE_NAME^}Usecase(repo domain.${SERVICE_NAME^}Repository) *${SERVICE_NAME^}Usecase {
	return &${SERVICE_NAME^}Usecase{repo: repo}
}

// Get${SERVICE_NAME^}ByID retrieves a ${SERVICE_NAME} by ID
func (uc *${SERVICE_NAME^}Usecase) Get${SERVICE_NAME^}ByID(id string) (*domain.${SERVICE_NAME^}, error) {
	return uc.repo.FindByID(id)
}

// Create${SERVICE_NAME^} creates a new ${SERVICE_NAME}
func (uc *${SERVICE_NAME^}Usecase) Create${SERVICE_NAME^}(${SERVICE_NAME} *domain.${SERVICE_NAME^}) error {
	return uc.repo.Store(${SERVICE_NAME})
}
EOF

echo "Created usecase file: $USECASE_FILE"

# Create application service
APPLICATION_FILE="internal/application/${SERVICE_NAME}_service.go"
cat > $APPLICATION_FILE << EOF
package application

import (
	"app-hexagonal/internal/domain"
	"app-hexagonal/internal/usecase"
)

// ${SERVICE_NAME^}Service provides application-level operations for ${SERVICE_NAME}s
type ${SERVICE_NAME^}Service struct {
	${SERVICE_NAME}Usecase usecase.${SERVICE_NAME^}UsecaseInterface
}

// New${SERVICE_NAME^}Service creates a new ${SERVICE_NAME} service
func New${SERVICE_NAME^}Service(${SERVICE_NAME}Usecase usecase.${SERVICE_NAME^}UsecaseInterface) *${SERVICE_NAME^}Service {
	return &${SERVICE_NAME^}Service{
		${SERVICE_NAME}Usecase: ${SERVICE_NAME}Usecase,
	}
}

// Get${SERVICE_NAME^}ByID retrieves a ${SERVICE_NAME} by their ID
func (s *${SERVICE_NAME^}Service) Get${SERVICE_NAME^}ByID(id string) (*domain.${SERVICE_NAME^}, error) {
	return s.${SERVICE_NAME}Usecase.Get${SERVICE_NAME^}ByID(id)
}

// Create${SERVICE_NAME^} creates a new ${SERVICE_NAME}
func (s *${SERVICE_NAME^}Service) Create${SERVICE_NAME^}(${SERVICE_NAME} *domain.${SERVICE_NAME^}) error {
	return s.${SERVICE_NAME}Usecase.Create${SERVICE_NAME^}(${SERVICE_NAME})
}
EOF

echo "Created application service file: $APPLICATION_FILE"

# Create gRPC service
GRPC_FILE="internal/delivery/grpc/${SERVICE_NAME}_service.go"
cat > $GRPC_FILE << EOF
package grpc

import (
	"context"

	v1 "app-hexagonal/api/v1"
	"app-hexagonal/internal/application"
	"app-hexagonal/internal/domain"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
)

// ${SERVICE_NAME^}ServiceServer implements the ${SERVICE_NAME^}Service gRPC service
type ${SERVICE_NAME^}ServiceServer struct {
	v1.Unimplemented${SERVICE_NAME^}ServiceServer
	${SERVICE_NAME}Service *application.${SERVICE_NAME^}Service
	logger      *zap.Logger
}

// New${SERVICE_NAME^}ServiceServer creates a new ${SERVICE_NAME^}ServiceServer
func New${SERVICE_NAME^}ServiceServer(${SERVICE_NAME}Service *application.${SERVICE_NAME^}Service, logger *zap.Logger) *${SERVICE_NAME^}ServiceServer {
	return &${SERVICE_NAME^}ServiceServer{
		${SERVICE_NAME}Service: ${SERVICE_NAME}Service,
		logger:      logger,
	}
}

// Get${SERVICE_NAME^} retrieves a ${SERVICE_NAME} by ID
func (s *${SERVICE_NAME^}ServiceServer) Get${SERVICE_NAME^}(ctx context.Context, req *v1.Get${SERVICE_NAME^}Request) (*v1.Get${SERVICE_NAME^}Response, error) {
	s.logger.Info("gRPC: Getting ${SERVICE_NAME} by ID", zap.String("${SERVICE_NAME}_id", req.GetId()))

	${SERVICE_NAME}, err := s.${SERVICE_NAME}Service.Get${SERVICE_NAME^}ByID(req.GetId())
	if err != nil {
		s.logger.Error("gRPC: Failed to get ${SERVICE_NAME}", zap.String("${SERVICE_NAME}_id", req.GetId()), zap.Error(err))
		return &v1.Get${SERVICE_NAME^}Response{
			Error:   true,
			Code:    int32(codes.NotFound),
			Message: "${SERVICE_NAME^} not found",
		}, nil
	}

	// Convert domain ${SERVICE_NAME} to protobuf ${SERVICE_NAME}
	proto${SERVICE_NAME^} := &v1.${SERVICE_NAME^}{
		Id:    ${SERVICE_NAME}.ID,
		// Add other fields here
	}

	s.logger.Info("gRPC: Successfully retrieved ${SERVICE_NAME}", zap.String("${SERVICE_NAME}_id", ${SERVICE_NAME}.ID))

	return &v1.Get${SERVICE_NAME^}Response{
		Error:   false,
		Code:    int32(codes.OK),
		Message: "${SERVICE_NAME^} retrieved successfully",
		Data:    proto${SERVICE_NAME^},
	}, nil
}

// Create${SERVICE_NAME^} creates a new ${SERVICE_NAME}
func (s *${SERVICE_NAME^}ServiceServer) Create${SERVICE_NAME^}(ctx context.Context, req *v1.Create${SERVICE_NAME^}Request) (*v1.Create${SERVICE_NAME^}Response, error) {
	s.logger.Info("gRPC: Creating new ${SERVICE_NAME}")

	// Create domain ${SERVICE_NAME}
	${SERVICE_NAME} := &domain.${SERVICE_NAME^}{
		ID:    "", // ID will be generated by the repository
		// Add other fields here
	}

	// Save ${SERVICE_NAME}
	err := s.${SERVICE_NAME}Service.Create${SERVICE_NAME^}(${SERVICE_NAME})
	if err != nil {
		s.logger.Error("gRPC: Failed to create ${SERVICE_NAME}", zap.Error(err))
		return &v1.Create${SERVICE_NAME^}Response{
			Error:   true,
			Code:    int32(codes.Internal),
			Message: "Failed to create ${SERVICE_NAME}",
		}, nil
	}

	// Convert domain ${SERVICE_NAME} to protobuf ${SERVICE_NAME}
	proto${SERVICE_NAME^} := &v1.${SERVICE_NAME^}{
		Id:    ${SERVICE_NAME}.ID,
		// Add other fields here
	}

	s.logger.Info("gRPC: ${SERVICE_NAME^} created successfully", zap.String("${SERVICE_NAME}_id", ${SERVICE_NAME}.ID))

	return &v1.Create${SERVICE_NAME^}Response{
		Error:   false,
		Code:    int32(codes.OK),
		Message: "${SERVICE_NAME^} created successfully",
		Data:    proto${SERVICE_NAME^},
	}, nil
}
EOF

echo "Created gRPC service file: $GRPC_FILE"

echo "API creation completed successfully!"
echo "Next steps:"
echo "1. Update the proto file with your specific message definitions"
echo "2. Update the domain model for ${SERVICE_NAME}"
echo "3. Update the repository interface and implementation"
echo "4. Register the service in the gRPC server"
echo "5. Update the main function to create and pass the new services"