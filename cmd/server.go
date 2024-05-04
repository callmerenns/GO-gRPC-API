package main

import (
	"log"
	"net"

	servicePb "github.com/altsaqif/go-grpc/pb/service"

	"github.com/altsaqif/go-grpc/cmd/configs"
	"github.com/altsaqif/go-grpc/cmd/middlewares"
	"github.com/altsaqif/go-grpc/cmd/services"
	"google.golang.org/grpc"
)

func main() {
	PORT := configs.GoDotEnvVariable("APP_PORT")

	netListen, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("Failed to Listened %v", err.Error())
	}

	db := configs.ConnectDatabase()

	// Implementasi Middleware
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middlewares.TokenInterceptor),
	)

	// Service Product
	serviceProduct := &services.ProductService{DB: db}
	servicePb.RegisterProductServiceServer(grpcServer, serviceProduct)

	// Service User
	serviceUser := &services.Server{DB: db}
	servicePb.RegisterAuthentificationServiceServer(grpcServer, serviceUser)

	log.Printf("Server Started at %v", netListen.Addr())

	if err := grpcServer.Serve(netListen); err != nil {
		log.Fatalf("Failed to Serve %v", err.Error())
	}
}
