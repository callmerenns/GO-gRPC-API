package main

import (
	"fmt"
	"log"
	"net"

	"github.com/altsaqif/go-grpc/cmd/configs"
	"github.com/altsaqif/go-grpc/cmd/delivery/middlewares"
	"github.com/altsaqif/go-grpc/cmd/entity"
	server "github.com/altsaqif/go-grpc/cmd/service"
	"github.com/altsaqif/go-grpc/cmd/shared/service"
	"github.com/altsaqif/go-grpc/cmd/utils"
	pb "github.com/altsaqif/go-grpc/pb/service"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type serviceServer struct {
	server.AuthServer
	server.ProductServer
	server.UserServer
	db              *gorm.DB
	jwtService      *service.JwtService
	rolePermissions map[string][]string
	tokenBlacklist  *utils.Blacklist
}

func NewServer(db *gorm.DB, jwtService *service.JwtService, rolePermissions map[string][]string, blacklist *utils.Blacklist) *serviceServer {
	return &serviceServer{
		AuthServer:      *server.NewAuthServer(db, jwtService, blacklist),
		ProductServer:   *server.NewProductServer(db),
		UserServer:      *server.NewUserServer(db),
		db:              db,
		jwtService:      jwtService,
		rolePermissions: rolePermissions,
		tokenBlacklist:  blacklist,
	}
}

func main() {
	// Load configuration
	cfg, err := configs.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set up database connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DbConfig.User, cfg.DbConfig.Password, cfg.DbConfig.Host, cfg.DbConfig.Port, cfg.DbConfig.Name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Drop existing tables
	db.Migrator().DropTable(&entity.User{}, &entity.Product{}, &entity.Enrollment{})

	err = db.AutoMigrate(&entity.User{}, &entity.Product{}, &entity.Enrollment{})
	if err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// Set up JWT service
	jwtService := service.NewJwtService(&service.TokenConfig{
		SecretKey: string(cfg.TokenConfig.JwtSignatureKey),
		Issuer:    cfg.TokenConfig.IssuerName,
	})

	// Define role permissions
	rolePermissions := map[string][]string{
		// Role Profiles
		"/go_grpc.UsersService/GetAllProfiles": {"admin"},
		"/go_grpc.UsersService/GetProfileById": {"admin"},

		// Role Products
		"/go_grpc.ProductsService/GetAllProducts":     {"admin", "costumer", "reseller"},
		"/go_grpc.ProductsService/GetProductById":     {"admin", "costumer", "reseller"},
		"/go_grpc.ProductsService/GetProductsByStock": {"admin", "reseller"},
		"/go_grpc.ProductsService/CreateProduct":      {"admin", "reseller"},
		"/go_grpc.ProductsService/UpdateProduct":      {"admin", "reseller"},
		"/go_grpc.ProductsService/DeleteProduct":      {"admin", "reseller"},
	}

	// Create blacklist
	blacklist := utils.NewBlacklist()

	// Create gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				middlewares.NewTokenInterceptor(string(cfg.TokenConfig.JwtSignatureKey), rolePermissions, blacklist),
				middlewares.AuthInterceptor(jwtService, db, rolePermissions),
			),
		),
	)
	pb.RegisterAuthenticationServiceServer(server, NewServer(db, jwtService, rolePermissions, blacklist))
	pb.RegisterProductsServiceServer(server, NewServer(db, jwtService, rolePermissions, blacklist))
	pb.RegisterUsersServiceServer(server, NewServer(db, jwtService, rolePermissions, blacklist))

	// Register reflection service on gRPC server
	reflection.Register(server)

	log.Println("gRPC server listening on :50051")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
