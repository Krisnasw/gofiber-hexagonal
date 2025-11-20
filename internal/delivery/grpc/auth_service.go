package grpc

import (
	"context"

	v1 "app-hexagonal/api/v1"
	"app-hexagonal/internal/application"
	"app-hexagonal/internal/domain"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
)

// AuthServiceServer implements the AuthService gRPC service
type AuthServiceServer struct {
	v1.UnimplementedAuthServiceServer
	authService *application.AuthService
	logger      *zap.Logger
}

// NewAuthServiceServer creates a new AuthServiceServer
func NewAuthServiceServer(authService *application.AuthService, logger *zap.Logger) *AuthServiceServer {
	return &AuthServiceServer{
		authService: authService,
		logger:      logger,
	}
}

// Login authenticates a user and returns tokens
func (s *AuthServiceServer) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error) {
	s.logger.Info("gRPC: Login request", zap.String("email", req.GetCredentials().GetEmail()))

	// Create credentials from request
	credentials := &domain.Credentials{
		Email:    req.GetCredentials().GetEmail(),
		Password: req.GetCredentials().GetPassword(),
	}

	// Authenticate user
	tokenResponse, err := s.authService.Login(credentials)
	if err != nil {
		s.logger.Error("gRPC: Login failed", zap.String("email", req.GetCredentials().GetEmail()), zap.Error(err))
		return &v1.LoginResponse{
			Error:   true,
			Code:    int32(codes.Unauthenticated),
			Message: "Invalid credentials",
		}, nil
	}

	// Convert to protobuf response
	tokenData := &v1.TokenData{
		AccessToken:  tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
		TokenType:    tokenResponse.TokenType,
		ExpiresIn:    int32(tokenResponse.ExpiresIn),
	}

	s.logger.Info("gRPC: Login successful", zap.String("email", req.GetCredentials().GetEmail()))

	return &v1.LoginResponse{
		Error:   false,
		Code:    int32(codes.OK),
		Message: "Login successful",
		Data:    tokenData,
	}, nil
}

// RefreshToken refreshes an access token using a refresh token
func (s *AuthServiceServer) RefreshToken(ctx context.Context, req *v1.RefreshTokenRequest) (*v1.RefreshTokenResponse, error) {
	s.logger.Info("gRPC: Refresh token request")

	// Refresh token
	tokenResponse, err := s.authService.RefreshToken(req.GetRefreshToken())
	if err != nil {
		s.logger.Error("gRPC: Token refresh failed", zap.Error(err))
		return &v1.RefreshTokenResponse{
			Error:   true,
			Code:    int32(codes.Unauthenticated),
			Message: "Invalid refresh token",
		}, nil
	}

	// Convert to protobuf response
	tokenData := &v1.TokenData{
		AccessToken:  tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
		TokenType:    tokenResponse.TokenType,
		ExpiresIn:    int32(tokenResponse.ExpiresIn),
	}

	s.logger.Info("gRPC: Token refresh successful")

	return &v1.RefreshTokenResponse{
		Error:   false,
		Code:    int32(codes.OK),
		Message: "Token refreshed successfully",
		Data:    tokenData,
	}, nil
}

// Logout invalidates the user's tokens
func (s *AuthServiceServer) Logout(ctx context.Context, req *v1.LogoutRequest) (*v1.LogoutResponse, error) {
	s.logger.Info("gRPC: Logout request")

	// Logout user
	err := s.authService.Logout(req.GetAccessToken())
	if err != nil {
		s.logger.Error("gRPC: Logout failed", zap.Error(err))
		return &v1.LogoutResponse{
			Error:   true,
			Code:    int32(codes.Unauthenticated),
			Message: "Invalid token",
		}, nil
	}

	s.logger.Info("gRPC: Logout successful")

	return &v1.LogoutResponse{
		Error:   false,
		Code:    int32(codes.OK),
		Message: "Logout successful",
	}, nil
}
