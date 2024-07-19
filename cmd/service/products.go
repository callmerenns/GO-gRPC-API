package service

import (
	"context"
	"time"

	"github.com/altsaqif/go-grpc/cmd/entity"
	pb "github.com/altsaqif/go-grpc/pb/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type ProductServer struct {
	pb.UnimplementedProductsServiceServer
	db *gorm.DB
}

func NewProductServer(db *gorm.DB) *ProductServer {
	return &ProductServer{
		db: db,
	}
}

// CreateProduct implementation
func (s *ProductServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductResponse, error) {
	// Ambil informasi pengguna dari konteks
	user, ok := ctx.Value("user").(*entity.User)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated: user not found in context")
	}

	// Buat produk baru
	product := entity.Product{
		Name:        req.Name,
		Description: req.Description,
		Stock:       req.Stock,
		Price:       req.Price,
		UserID:      user.ID, // Tambahkan ID pengguna yang membuat produk
	}

	// Gunakan goroutine untuk menyimpan produk dan membuat entri di tabel pivot enrollments secara paralel
	errChan := make(chan error, 1)
	productChan := make(chan entity.Product, 1)

	go func() {
		if err := s.db.Create(&product).Error; err != nil {
			errChan <- err
			return
		}
		productChan <- product
		errChan <- nil
	}()

	// Periksa apakah ada error dari goroutine
	if err := <-errChan; err != nil {
		return nil, err
	}

	// Ambil hasil produk dari goroutine
	product = <-productChan

	// Buat entri di tabel pivot enrollments
	enrollment := entity.Enrollment{
		UserID:    user.ID,
		ProductID: product.ID,
	}

	// Gunakan goroutine untuk menyimpan entri enrollment
	go func() {
		if err := s.db.Create(&enrollment).Error; err != nil {
			errChan <- err
			return
		}
		errChan <- nil
	}()

	// Periksa apakah ada error dari goroutine
	if err := <-errChan; err != nil {
		return nil, err
	}

	return &pb.ProductResponse{
		Product: &pb.Product{
			Id:          uint64(product.ID),
			Name:        product.Name,
			Description: product.Description,
			Stock:       product.Stock,
			Price:       product.Price,
			UserId:      uint64(user.ID),
			CreatedAt:   ToProtoTimestamp(product.CreatedAt),
			UpdatedAt:   ToProtoTimestamp(product.UpdatedAt),
		},
	}, nil
}

// GetAllProducts implementation
func (s *ProductServer) GetAllProducts(ctx context.Context, req *pb.GetAllProductsRequest) (*pb.GetAllProductsResponse, error) {
	var products []entity.Product
	var total int64

	// Channel untuk mengumpulkan error dari goroutine
	errChan := make(chan error, 2)

	// Hitung total produk
	go func() {
		if err := s.db.Model(&entity.Product{}).Count(&total).Error; err != nil {
			errChan <- err
			return
		}
		errChan <- nil
	}()

	// Ambil produk dengan limit dan offset serta preload data pengguna melalui tabel pivot
	go func() {
		if err := s.db.Preload("Users").Limit(int(req.Limit)).Offset(int(req.Offset)).Find(&products).Error; err != nil {
			errChan <- err
			return
		}
		errChan <- nil
	}()

	// Tunggu sampai kedua goroutine selesai
	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			return nil, err
		}
	}

	var productResponses []*pb.Product
	for _, product := range products {
		var userProfiles []*pb.UserProfile
		for _, user := range product.Users {
			userProfiles = append(userProfiles, &pb.UserProfile{
				Id:        uint64(user.ID),
				Firstname: user.FirstName,
				Lastname:  user.LastName,
				Email:     user.Email,
				Role:      user.Role,
				CreatedAt: ToProtoTimestamp(user.CreatedAt),
				UpdatedAt: ToProtoTimestamp(user.UpdatedAt),
			})
		}
		productResponses = append(productResponses, &pb.Product{
			Id:          uint64(product.ID),
			Name:        product.Name,
			Description: product.Description,
			Stock:       product.Stock,
			Price:       product.Price,
			UserId:      uint64(product.UserID),
			Users:       userProfiles,
			CreatedAt:   ToProtoTimestamp(product.CreatedAt),
			UpdatedAt:   ToProtoTimestamp(product.UpdatedAt),
		})
	}

	return &pb.GetAllProductsResponse{
		Products: productResponses,
		Total:    int32(total),
		Limit:    req.Limit,
		Offset:   req.Offset,
	}, nil
}

// GetProductById implementation
func (s *ProductServer) GetProductById(ctx context.Context, req *pb.GetProductByIdRequest) (*pb.ProductResponse, error) {
	var product entity.Product

	// Channel untuk mengumpulkan error dari goroutine
	errChan := make(chan error, 1)

	// Temukan produk berdasarkan ID dan preload data pengguna terkait
	go func() {
		if err := s.db.Preload("Users").First(&product, req.Id).Error; err != nil {
			errChan <- err
			return
		}
		errChan <- nil
	}()

	// Tunggu goroutine selesai
	if err := <-errChan; err != nil {
		return nil, err
	}

	var userProfiles []*pb.UserProfile
	for _, user := range product.Users {
		userProfiles = append(userProfiles, &pb.UserProfile{
			Id:        uint64(user.ID),
			Firstname: user.FirstName,
			Lastname:  user.LastName,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: ToProtoTimestamp(user.CreatedAt),
			UpdatedAt: ToProtoTimestamp(user.UpdatedAt),
		})
	}

	return &pb.ProductResponse{
		Product: &pb.Product{
			Id:          uint64(product.ID),
			Name:        product.Name,
			Description: product.Description,
			Stock:       product.Stock,
			Price:       product.Price,
			CreatedAt:   ToProtoTimestamp(product.CreatedAt),
			UpdatedAt:   ToProtoTimestamp(product.UpdatedAt),
			Users:       userProfiles,
		},
	}, nil
}

// GetProductsByStock implementation
func (s *ProductServer) GetProductsByStock(ctx context.Context, req *pb.GetProductsByStockRequest) (*pb.GetProductsByStockResponse, error) {
	var products []entity.Product

	// Channel untuk mengumpulkan error dari goroutine
	errChan := make(chan error, 1)

	// Goroutine untuk menemukan produk berdasarkan stok dan preload data pengguna terkait
	go func() {
		if err := s.db.Where("stock = ?", req.Stock).Preload("Users").Find(&products).Error; err != nil {
			errChan <- err
			return
		}
		errChan <- nil
	}()

	// Tunggu goroutine selesai dan periksa error
	if err := <-errChan; err != nil {
		return nil, err
	}

	// Memproses hasil dari goroutine
	var productResponses []*pb.Product
	for _, product := range products {
		var userProfiles []*pb.UserProfile
		for _, user := range product.Users {
			userProfiles = append(userProfiles, &pb.UserProfile{
				Id:        uint64(user.ID),
				Firstname: user.FirstName,
				Lastname:  user.LastName,
				Email:     user.Email,
				Role:      user.Role,
				CreatedAt: ToProtoTimestamp(user.CreatedAt),
				UpdatedAt: ToProtoTimestamp(user.UpdatedAt),
			})
		}

		productResponses = append(productResponses, &pb.Product{
			Id:          uint64(product.ID),
			Name:        product.Name,
			Description: product.Description,
			Stock:       product.Stock,
			Price:       product.Price,
			CreatedAt:   ToProtoTimestamp(product.CreatedAt),
			UpdatedAt:   ToProtoTimestamp(product.UpdatedAt),
			Users:       userProfiles,
		})
	}

	return &pb.GetProductsByStockResponse{
		Products: productResponses,
	}, nil
}

func (s *ProductServer) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.ProductResponse, error) {
	var product entity.Product

	// Channel untuk mengumpulkan hasil dari goroutine
	resultChan := make(chan error, 1)

	// Goroutine untuk mencari produk berdasarkan ID
	go func() {
		if err := s.db.First(&product, req.Id).Error; err != nil {
			resultChan <- status.Errorf(codes.NotFound, "Product not found")
			return
		}
		resultChan <- nil
	}()

	// Tunggu goroutine selesai dan periksa error
	if err := <-resultChan; err != nil {
		return nil, err
	}

	// Update informasi produk
	product.Name = req.Name
	product.Description = req.Description
	product.Stock = req.Stock
	product.Price = req.Price

	// Goroutine untuk menyimpan perubahan produk
	go func() {
		if err := s.db.Save(&product).Error; err != nil {
			resultChan <- status.Errorf(codes.Internal, "Failed to update product: %v", err)
			return
		}
		resultChan <- nil
	}()

	// Tunggu goroutine selesai dan periksa error
	if err := <-resultChan; err != nil {
		return nil, err
	}

	return &pb.ProductResponse{
		Product: &pb.Product{
			Id:          uint64(product.ID),
			Name:        product.Name,
			Description: product.Description,
			Stock:       product.Stock,
			Price:       product.Price,
			CreatedAt:   ToProtoTimestamp(product.CreatedAt),
			UpdatedAt:   ToProtoTimestamp(product.UpdatedAt),
		},
	}, nil
}

func (s *ProductServer) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	var product entity.Product

	// Channel untuk mengumpulkan hasil dari goroutine
	resultChan := make(chan error, 1)

	// Goroutine untuk mencari produk berdasarkan ID
	go func() {
		if err := s.db.First(&product, req.Id).Error; err != nil {
			resultChan <- status.Errorf(codes.NotFound, "Product not found")
			return
		}
		resultChan <- nil
	}()

	// Tunggu goroutine selesai dan periksa error
	if err := <-resultChan; err != nil {
		return nil, err
	}

	// Goroutine untuk menghapus produk
	go func() {
		if err := s.db.Delete(&product).Error; err != nil {
			resultChan <- status.Errorf(codes.Internal, "Failed to delete product: %v", err)
			return
		}
		resultChan <- nil
	}()

	// Tunggu goroutine selesai dan periksa error
	if err := <-resultChan; err != nil {
		return nil, err
	}

	return &pb.DeleteProductResponse{
		Message: "Product successfully deleted",
	}, nil
}

// Convert time.Time to *timestamppb.Timestamp
func ToProtoTimestamp(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}

// Convert *timestamppb.Timestamp to time.Time
func FromProtoTimestamp(t *timestamppb.Timestamp) time.Time {
	return t.AsTime()
}
