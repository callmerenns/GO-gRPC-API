package middlewares

import (
	"context"
	"log"
	"strings"

	"github.com/altsaqif/go-grpc/cmd/entity"
	"github.com/altsaqif/go-grpc/cmd/shared/service"
	"github.com/altsaqif/go-grpc/cmd/utils"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type contextKey string

const (
	userClaimsKey contextKey = "user"
)

func NewTokenInterceptor(secretKey string, rolePermissions map[string][]string, blacklist *utils.Blacklist) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		allowedServices := map[string]bool{
			"/go_grpc.AuthenticationService/Login":    true,
			"/go_grpc.AuthenticationService/Register": true,
			"/go_grpc.AuthenticationService/Logout":   true,
		}

		if allowedServices[info.FullMethod] {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			log.Println("No metadata found in context")
			return nil, status.Errorf(codes.Unauthenticated, "no metadata found in context")
		}

		authHeader, ok := md["authorization"]
		if !ok || len(authHeader) == 0 {
			log.Println("No authorization token provided")
			return nil, status.Errorf(codes.Unauthenticated, "no authorization token provided")
		}

		log.Printf("Authorization header: %v", authHeader[0])

		parts := strings.Split(authHeader[0], " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			log.Println("Invalid authorization header format")
			return nil, status.Errorf(codes.Unauthenticated, "invalid authorization header format")
		}

		token := parts[1]

		if blacklist.Exists(token) {
			log.Printf("Token is blacklisted: %v", token)
			return nil, status.Errorf(codes.Unauthenticated, "token is blacklisted")
		}

		tokenObj, err := jwt.ParseWithClaims(token, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
		if err != nil {
			log.Printf("Invalid token: %v", err)
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		if claims, ok := tokenObj.Claims.(jwt.MapClaims); ok && tokenObj.Valid {
			role, roleExists := claims["role"].(string)
			if !roleExists {
				log.Println("Role claim not found in token")
				return nil, status.Errorf(codes.Unauthenticated, "role claim not found in token")
			}

			if !isRoleAllowed(role, info.FullMethod, rolePermissions) {
				log.Printf("Role %s is not allowed to access %s", role, info.FullMethod)
				return nil, status.Errorf(codes.PermissionDenied, "role %s is not allowed to access %s", role, info.FullMethod)
			}

			ctx = context.WithValue(ctx, userClaimsKey, claims)
			log.Printf("Token claims set in context: %v", claims)
		} else {
			log.Println("Invalid token claims")
			return nil, status.Errorf(codes.Unauthenticated, "invalid token claims")
		}

		return handler(ctx, req)
	}
}

func AuthInterceptor(jwtService *service.JwtService, db *gorm.DB, rolePermissions map[string][]string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		method := info.FullMethod
		if shouldSkipAuth(method) {
			return handler(ctx, req)
		}

		token, err := getTokenFromContext(ctx)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "unauthenticated: %v", err)
		}

		claims, err := jwtService.ParseToken(token)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		role, roleExists := claims["role"].(string)
		if !roleExists {
			return nil, status.Errorf(codes.Unauthenticated, "role claim not found in token")
		}

		if !isRoleAllowed(role, method, rolePermissions) {
			return nil, status.Errorf(codes.PermissionDenied, "role %s is not allowed to access %s", role, method)
		}

		// Ambil user ID dari klaim token
		userID, userIDExists := claims["userId"].(float64)
		if !userIDExists {
			return nil, status.Errorf(codes.Unauthenticated, "user_id claim not found in token")
		}

		// Ambil informasi pengguna dari database
		var user entity.User
		if err := db.Where("id = ?", uint(userID)).First(&user).Error; err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "user not found")
		}

		// Tambahkan informasi pengguna ke dalam konteks
		ctx = context.WithValue(ctx, "user", &user)
		log.Printf("User set in context: %v", user)

		return handler(ctx, req)
	}
}

func isRoleAllowed(role, method string, rolePermissions map[string][]string) bool {
	allowedRoles, methodExists := rolePermissions[method]
	if !methodExists {
		return false
	}

	for _, allowedRole := range allowedRoles {
		if role == allowedRole {
			return true
		}
	}

	return false
}

func shouldSkipAuth(method string) bool {
	methodsToSkip := []string{
		"/go_grpc.AuthenticationService/Login",
		"/go_grpc.AuthenticationService/Register",
		"/go_grpc.AuthenticationService/Logout",
	}

	for _, m := range methodsToSkip {
		if strings.HasPrefix(method, m) {
			return true
		}
	}
	return false
}

func getTokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "no metadata found in context")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return "", status.Errorf(codes.Unauthenticated, "no authorization header found")
	}

	parts := strings.Split(authHeader[0], " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", status.Errorf(codes.Unauthenticated, "invalid authorization header format")
	}

	return parts[1], nil
}
