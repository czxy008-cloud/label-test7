package service

import (
	"context"
	"fresh-groupbuy/internal/models"
	"fresh-groupbuy/internal/repository"
	"fresh-groupbuy/pkg/database"
)

type ProductService struct {
	productRepo *repository.ProductRepo
}

func NewProductService() *ProductService {
	return &ProductService{
		productRepo: repository.NewProductRepo(),
	}
}

func (s *ProductService) GetByID(ctx context.Context, tx database.DBTX, id int64) (*models.Product, error) {
	return s.productRepo.GetByID(ctx, tx, id)
}

func (s *ProductService) List(ctx context.Context, tx database.DBTX, category string, page, pageSize int) ([]*models.Product, int64, error) {
	offset := (page - 1) * pageSize
	return s.productRepo.List(ctx, tx, category, offset, pageSize)
}
