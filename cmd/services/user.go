package services

import (
	"context"
	"time"

	"github.com/altsaqif/go-grpc/cmd/helpers"
	"github.com/altsaqif/go-grpc/cmd/models"
	"github.com/altsaqif/go-grpc/pb/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type Server struct {
	service.UnimplementedAuthentificationServiceServer
	DB *gorm.DB
}

func (s *Server) Login(ctx context.Context, req *service.LoginRequest) (*service.LoginResponse, error) {
	var Users models.User
	if err := s.DB.Where("email = ?", req.Email).First(&Users).Error; err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid Email")
	}

	if !helpers.CheckPasswordHash(req.Password, Users.Password) {
		return nil, status.Error(codes.Unauthenticated, "Invalid Password")
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": 123,                                  // Example user ID
		"exp":     time.Now().Add(time.Hour * 1).Unix(), // Token expiration time
	})
	tokenString, err := token.SignedString([]byte("your-secret-key")) // Change "your-secret-key" to your actual secret key
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to generate token: %v", err)
	}

	// Store token in metadata
	md := metadata.Pairs("token", tokenString)
	grpc.SendHeader(ctx, md)

	// Return token in LoginResponse
	return &service.LoginResponse{
		Status: &service.Status{Status: 200, Message: "Success"},
		Token:  tokenString,
	}, nil
}

func (s *Server) Register(ctx context.Context, req *service.RegisterRequest) (*service.Status, error) {
	if req.Password != req.ConfirmPassword {
		return nil, status.Error(codes.InvalidArgument, "passwords do not match")
	}

	hashedPassword, err := helpers.HashPassword(req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	Users := models.User{
		ID:       uuid.New().String(),
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := s.DB.Create(&Users).Error; err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	var response service.Status
	response.Status = 1
	response.Message = "User registered successfully"

	return &response, nil
}

func (s *Server) Logout(ctx context.Context, req *service.LogoutRequest) (*service.Status, error) {
	// Retrieve token from metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Metadata is missing")
	}

	// Remove token from metadata
	delete(md, "token")

	// Call grpc.SetHeader to update metadata with the modified value
	grpc.SetHeader(ctx, md)

	// Return success status
	return &service.Status{Status: 200, Message: "User logged out successfully"}, nil
}
