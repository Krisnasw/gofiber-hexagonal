package grpc

import (
	"net"

	v1 "app-hexagonal/api/v1"
	"app-hexagonal/internal/application"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server represents the gRPC server
type Server struct {
	server *grpc.Server
	logger *zap.Logger
	port   string
}

// NewServer creates a new gRPC server
func NewServer(logger *zap.Logger, port string) *Server {
	return &Server{
		logger: logger,
		port:   port,
	}
}

// Start starts the gRPC server
func (s *Server) Start(userService *application.UserService, authService *application.AuthService) error {
	// Create a new gRPC server
	s.server = grpc.NewServer()

	// Register the user service
	userServiceServer := NewUserServiceServer(userService, s.logger)
	v1.RegisterUserServiceServer(s.server, userServiceServer)

	// Register the auth service
	authServiceServer := NewAuthServiceServer(authService, s.logger)
	v1.RegisterAuthServiceServer(s.server, authServiceServer)

	// Enable reflection for debugging
	reflection.Register(s.server)

	// Listen on the specified port
	lis, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		s.logger.Error("Failed to listen", zap.Error(err))
		return err
	}

	s.logger.Info("Starting gRPC server", zap.String("port", s.port))

	// Start serving
	if err := s.server.Serve(lis); err != nil {
		s.logger.Error("Failed to serve", zap.Error(err))
		return err
	}

	return nil
}

// Stop stops the gRPC server
func (s *Server) Stop() {
	s.logger.Info("Stopping gRPC server")
	s.server.GracefulStop()
}
