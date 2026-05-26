package service

import (
	"context"
	"fresh-groupbuy/internal/models"
	"fresh-groupbuy/internal/repository"
	"fresh-groupbuy/pkg/database"
)

type CartService struct {
	cartRepo    *repository.CartRepo
	productRepo *repository.ProductRepo
}

func NewCartService() *CartService {
	return &CartService{
		cartRepo:    repository.NewCartRepo(),
		productRepo: repository.NewProductRepo(),
	}
}

func (s *CartService) Add(ctx context.Context, tx database.DBTX, userID, productID int64, quantity int) error {
	if tx == nil {
		tx = database.DB
	}
	if quantity <= 0 {
		quantity = 1
	}

	product, err := s.productRepo.GetByID(ctx, tx, productID)
	if err != nil {
		return err
	}

	if product.Stock < quantity {
		return ErrInsufficientStock
	}

	return s.cartRepo.Add(ctx, tx, userID, productID, quantity)
}

func (s *CartService) UpdateQuantity(ctx context.Context, tx database.DBTX, userID, productID int64, quantity int) error {
	if tx == nil {
		tx = database.DB
	}
	if quantity <= 0 {
		return s.cartRepo.Remove(ctx, tx, userID, productID)
	}

	product, err := s.productRepo.GetByID(ctx, tx, productID)
	if err != nil {
		return err
	}

	if product.Stock < quantity {
		return ErrInsufficientStock
	}

	return s.cartRepo.UpdateQuantity(ctx, tx, userID, productID, quantity)
}

func (s *CartService) Remove(ctx context.Context, tx database.DBTX, userID, productID int64) error {
	if tx == nil {
		tx = database.DB
	}
	return s.cartRepo.Remove(ctx, tx, userID, productID)
}

func (s *CartService) List(ctx context.Context, tx database.DBTX, userID int64) ([]*models.Cart, error) {
	if tx == nil {
		tx = database.DB
	}
	return s.cartRepo.List(ctx, tx, userID)
}

func (s *CartService) Clear(ctx context.Context, tx database.DBTX, userID int64) error {
	if tx == nil {
		tx = database.DB
	}
	return s.cartRepo.Clear(ctx, tx, userID)
}
