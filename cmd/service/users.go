package service

import (
	"context"
	"log"

	"github.com/altsaqif/go-grpc/cmd/entity"
	pb "github.com/altsaqif/go-grpc/pb/service"
	"gorm.io/gorm"
)

type UserServer struct {
	pb.UnimplementedUsersServiceServer
	db *gorm.DB
}

func NewUserServer(db *gorm.DB) *UserServer {
	return &UserServer{
		db: db,
	}
}

// GetAllProfiles implementation
func (s *UserServer) GetAllProfiles(ctx context.Context, req *pb.GetAllProfilesRequest) (*pb.GetAllProfilesResponse, error) {
	var users []entity.User
	var total int64
	errChan := make(chan error, 1)

	// Hitung total pengguna
	go func() {
		if err := s.db.Model(&entity.User{}).Count(&total).Error; err != nil {
			errChan <- err
			return
		}
		errChan <- nil
	}()

	// Ambil pengguna dengan limit dan offset serta preload data produk mereka
	go func() {
		if err := s.db.Preload("Products").Limit(int(req.Limit)).Offset(int(req.Offset)).Find(&users).Error; err != nil {
			errChan <- err
			return
		}
		errChan <- nil
	}()

	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			return nil, err
		}
	}

	log.Printf("Total users: %d\n", total)
	log.Printf("Users fetched: %v\n", users)

	var userProfiles []*pb.UserProfile
	profileChan := make(chan *pb.UserProfile, len(users))

	for _, user := range users {
		go func(user entity.User) {
			var products []*pb.Product
			for _, product := range user.Products {
				products = append(products, &pb.Product{
					Id:          uint64(product.ID),
					Name:        product.Name,
					Description: product.Description,
					Stock:       product.Stock,
					Price:       product.Price,
					UserId:      uint64(product.UserID),
					CreatedAt:   ToProtoTimestamp(product.CreatedAt),
					UpdatedAt:   ToProtoTimestamp(product.UpdatedAt),
				})
			}

			profileChan <- &pb.UserProfile{
				Id:        uint64(user.ID),
				Firstname: user.FirstName,
				Lastname:  user.LastName,
				Email:     user.Email,
				Role:      user.Role,
				Products:  products,
				CreatedAt: ToProtoTimestamp(user.CreatedAt),
				UpdatedAt: ToProtoTimestamp(user.UpdatedAt),
			}
		}(user)
	}

	for range users {
		userProfiles = append(userProfiles, <-profileChan)
	}

	return &pb.GetAllProfilesResponse{
		Users:  userProfiles,
		Total:  int32(total),
		Limit:  req.Limit,
		Offset: req.Offset,
	}, nil
}

// GetProfileById implementation
func (s *UserServer) GetProfileById(ctx context.Context, req *pb.GetProfileByIdRequest) (*pb.UserProfile, error) {
	var user entity.User
	errChan := make(chan error, 1)
	userChan := make(chan entity.User, 1)

	// Ambil pengguna dengan preload data produk mereka
	go func() {
		if err := s.db.Preload("Products").First(&user, req.Id).Error; err != nil {
			errChan <- err
			return
		}
		userChan <- user
		errChan <- nil
	}()

	// Periksa apakah ada error dari goroutine
	if err := <-errChan; err != nil {
		return nil, err
	}

	// Ambil hasil dari goroutine
	user = <-userChan

	// Buat slice products menggunakan goroutine
	productsChan := make(chan []*pb.Product, 1)

	go func() {
		var products []*pb.Product
		for _, product := range user.Products {
			products = append(products, &pb.Product{
				Id:          uint64(product.ID),
				Name:        product.Name,
				Description: product.Description,
				Stock:       product.Stock,
				Price:       product.Price,
				UserId:      uint64(product.UserID),
				CreatedAt:   ToProtoTimestamp(product.CreatedAt),
				UpdatedAt:   ToProtoTimestamp(product.UpdatedAt),
			})
		}
		productsChan <- products
	}()

	products := <-productsChan

	return &pb.UserProfile{
		Id:        uint64(user.ID),
		Firstname: user.FirstName,
		Lastname:  user.LastName,
		Email:     user.Email,
		Role:      user.Role,
		Products:  products,
		CreatedAt: ToProtoTimestamp(user.CreatedAt),
		UpdatedAt: ToProtoTimestamp(user.UpdatedAt),
	}, nil
}
