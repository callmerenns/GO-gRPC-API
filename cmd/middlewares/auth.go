package middlewares

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TokenInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Daftar service yang diizinkan tanpa token
	allowedServices := map[string]bool{
		"/go_grpc.AuthentificationService/Login":    true,
		"/go_grpc.AuthentificationService/Register": true,
		"/go_grpc.AuthentificationService/Logout":   true,
	}

	// Cek jika service diizinkan tanpa token
	if allowedServices[info.FullMethod] {
		return handler(ctx, req)
	}

	// Extract token from context metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Metadata is missing")
	}
	tokenString := md.Get("token")
	if len(tokenString) == 0 {
		return nil, status.Error(codes.Unauthenticated, "Token is missing")
	}

	// Add token to context for subsequent handlers to use
	ctx = context.WithValue(ctx, "token", tokenString[0])

	// Call the handler
	return handler(ctx, req)
}
