package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/altsaqif/go-grpc/cmd/entity"
	"github.com/altsaqif/go-grpc/cmd/shared/service"
	"github.com/altsaqif/go-grpc/cmd/utils"
	pb "github.com/altsaqif/go-grpc/pb/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type AuthServer struct {
	pb.UnimplementedAuthenticationServiceServer
	db             *gorm.DB
	jwtService     service.JwtService
	tokenBlacklist utils.Blacklist
}

func NewAuthServer(db *gorm.DB, jwtService *service.JwtService, blacklist *utils.Blacklist) *AuthServer {
	return &AuthServer{
		db:             db,
		jwtService:     *jwtService,
		tokenBlacklist: *blacklist,
	}
}

// Login implementation
func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	type result struct {
		user entity.User
		err  error
	}

	resCh := make(chan result)

	// Menjalankan query database dalam goroutine
	go func() {
		var user entity.User
		if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
			resCh <- result{err: status.Errorf(codes.NotFound, "user not found: %v", err)}
			return
		}

		if !utils.CheckPasswordHash(req.Password, user.Password) {
			resCh <- result{err: status.Errorf(codes.Unauthenticated, "invalid credentials")}
			return
		}

		resCh <- result{user: user}
	}()

	// Menunggu hasil dari goroutine
	res := <-resCh
	if res.err != nil {
		return nil, res.err
	}

	user := res.user

	// Membuat token dalam goroutine
	tokenCh := make(chan struct {
		token string
		err   error
	})

	go func() {
		token, err := s.jwtService.GenerateToken(user.ID, user.Role)
		tokenCh <- struct {
			token string
			err   error
		}{token, err}
	}()

	// Menunggu hasil dari pembuatan token
	tokenRes := <-tokenCh
	if tokenRes.err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token: %v", tokenRes.err)
	}

	token := tokenRes.token

	// Set token in response metadata
	md := metadata.Pairs("authorization", "Bearer "+token)
	grpc.SendHeader(ctx, md)

	return &pb.LoginResponse{
		Token: token,
		User: &pb.UserProfile{
			Id:        uint64(user.ID),
			Firstname: user.FirstName,
			Lastname:  user.LastName,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: ToProtoTimestamp(user.CreatedAt),
			UpdatedAt: ToProtoTimestamp(user.UpdatedAt),
		},
	}, nil
}

// Register implementation
func (s *AuthServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// Memastikan password dan konfirmasi password cocok
	if req.Password != req.ConfirmPassword {
		return nil, fmt.Errorf("password and confirm password do not match")
	}

	type hashResult struct {
		hashedPassword string
		err            error
	}

	hashCh := make(chan hashResult)

	// Menjalankan proses hashing password dalam goroutine
	go func() {
		hashedPassword, err := utils.HashPassword(req.Password)
		hashCh <- hashResult{hashedPassword, err}
	}()

	// Menunggu hasil dari hashing password
	hashRes := <-hashCh
	if hashRes.err != nil {
		return nil, hashRes.err
	}

	user := entity.User{
		FirstName: req.Firstname,
		LastName:  req.Lastname,
		Email:     req.Email,
		Password:  hashRes.hashedPassword,
		Role:      req.Role,
	}

	userCh := make(chan error)

	// Menjalankan proses pembuatan pengguna dalam goroutine
	go func() {
		userCh <- s.db.Create(&user).Error
	}()

	// Menunggu hasil dari pembuatan pengguna
	if err := <-userCh; err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{
		User: &pb.UserProfile{
			Id:        uint64(user.ID),
			Firstname: user.FirstName,
			Lastname:  user.LastName,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: ToProtoTimestamp(user.CreatedAt),
			UpdatedAt: ToProtoTimestamp(user.UpdatedAt),
		},
	}, nil
}

// Logout implementation
func (s *AuthServer) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "no metadata found in context")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "no authorization token provided")
	}

	parts := strings.Split(authHeader[0], " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return nil, status.Errorf(codes.Unauthenticated, "invalid authorization header format")
	}

	token := parts[1]

	// Membuat channel untuk menunggu hasil penambahan token ke blacklist
	blacklistCh := make(chan struct {
		success bool
		err     error
	})

	// Menjalankan penambahan token ke dalam daftar blacklist dalam goroutine
	go func() {
		err := s.tokenBlacklist.Add(token)
		if err != nil {
			blacklistCh <- struct {
				success bool
				err     error
			}{false, err}
			return
		}
		blacklistCh <- struct {
			success bool
			err     error
		}{true, nil}
	}()

	// Menunggu hasil dari goroutine
	res := <-blacklistCh
	if !res.success {
		return nil, status.Errorf(codes.Internal, "failed to blacklist token: %v", res.err)
	}

	return &pb.LogoutResponse{Message: "Logout successful"}, nil
}
