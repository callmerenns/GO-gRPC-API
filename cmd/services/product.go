package services

import (
	"context"
	"log"

	pagingPb "github.com/altsaqif/go-grpc/pb/pagination"
	servicePb "github.com/altsaqif/go-grpc/pb/service"

	"github.com/altsaqif/go-grpc/cmd/helpers"
	"github.com/altsaqif/go-grpc/cmd/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type ProductService struct {
	servicePb.UnimplementedProductServiceServer
	DB *gorm.DB
}

// Get All Product
func (p *ProductService) GetProducts(ctx context.Context, pageParam *servicePb.Page) (*servicePb.Products, error) {

	var (
		page       int64 = 1
		pagination pagingPb.Pagination
		products   []*servicePb.Product
	)

	if pageParam.GetPage() != 0 {
		page = pageParam.GetPage()
	}

	sql := p.DB.Table("products AS p").
		Joins("LEFT JOIN categories AS c ON c.id = p.category_id").
		Select("p.id", "p.name", "p.price", "p.stock", "c.id as category_id", "c.name as category_name")

	offset, limit := helpers.Pagination(sql, page, &pagination)
	rows, err := sql.Offset(int(offset)).Limit(int(limit)).Rows()

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var (
			product  servicePb.Product
			category servicePb.Category
		)

		if err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.Stock, &category.Id, &category.Name); err != nil {
			log.Fatalf("Gagal Mengambil Row data %v", err.Error())
		}

		product.Category = &category
		products = append(products, &product)

	}
	response := &servicePb.Products{
		Pagination: &pagination,
		Data:       products,
	}

	return response, nil
}

func (p *ProductService) GetProduct(ctx context.Context, id *servicePb.Id) (*servicePb.Product, error) {
	row := p.DB.Table("products AS p").
		Joins("LEFT JOIN categories AS c ON c.id = p.category_id").
		Select("p.id", "p.name", "p.price", "p.stock", "c.id as category_id", "c.name as category_name").
		Where("p.id = ?", id.GetId()).
		Row()

	var (
		product  servicePb.Product
		category servicePb.Category
	)

	if err := row.Scan(&product.Id, &product.Name, &product.Price, &product.Stock, &category.Id, &category.Name); err != nil {
		return nil, status.Error(codes.NotFound, "Data Not Found In Database")
	}

	product.Category = &category

	return &product, nil
}

func (p *ProductService) CreateProduct(ctx context.Context, productData *servicePb.Product) (*servicePb.CProduct, error) {
	var (
		products   servicePb.CProduct
		product2   servicePb.Product
		categories servicePb.Category
	)

	err := p.DB.Transaction(func(tx *gorm.DB) error {
		category := servicePb.Category{
			Id:   0,
			Name: productData.GetCategory().GetName(),
		}

		if err := tx.Table("categories").
			Where("LCASE(name) = ?", category.GetName()).
			FirstOrCreate(&category).Error; err != nil {
			return err
		}

		product := models.Product{
			Id:          productData.GetId(),
			Name:        productData.GetName(),
			Stock:       productData.GetStock(),
			Price:       productData.GetPrice(),
			Category_id: category.GetId(),
		}

		if err := tx.Table("products").Create(&product).Error; err != nil {
			return err
		}

		// Get Data Category
		categories.Id = category.Id
		categories.Name = category.Name

		// Get Data Product
		product2.Id = product.Id
		product2.Name = product.Name
		product2.Stock = product.Stock
		product2.Price = product.Price
		product2.Category = &categories

		products.Status = 1
		products.Message = "Successfully Create Data"
		products.Product = &product2

		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &products, nil
}

func (p *ProductService) UpdateProduct(ctx context.Context, productData *servicePb.Product) (*servicePb.UProduct, error) {
	var (
		products   servicePb.UProduct
		product2   servicePb.Product
		categories servicePb.Category
	)

	err := p.DB.Transaction(func(tx *gorm.DB) error {
		category := servicePb.Category{
			Id:   0,
			Name: productData.GetCategory().GetName(),
		}

		if err := tx.Table("categories").
			Where("LCASE(name) = ?", category.GetName()).
			FirstOrCreate(&category).Error; err != nil {
			return err
		}

		product := models.Product{
			Id:          productData.GetId(),
			Name:        productData.GetName(),
			Stock:       productData.GetStock(),
			Price:       productData.GetPrice(),
			Category_id: category.GetId(),
		}

		if err := tx.Table("products").Where("id = ?", product.Id).Updates(&product).Error; err != nil {
			return status.Error(codes.NotFound, "Data Not Found In Database")
		}

		// Get Data Category
		categories.Id = category.Id
		categories.Name = category.Name

		// Get Data Product
		product2.Id = product.Id
		product2.Name = product.Name
		product2.Stock = product.Stock
		product2.Price = product.Price
		product2.Category = &categories

		products.Status = 1
		products.Message = "Successfully Update Data"
		products.Product = &product2

		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &products, nil
}

func (p *ProductService) DeleteProduct(ctx context.Context, id *servicePb.Id) (*servicePb.Status, error) {
	var response servicePb.Status

	if err := p.DB.Table("products").Where("id = ?", id.GetId()).Delete(nil).Error; err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	response.Status = 1
	response.Message = "Delete Data Successfully"
	return &response, nil
}
