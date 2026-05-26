package service

import (
	"context"
	"fresh-groupbuy/internal/models"
	"fresh-groupbuy/internal/repository"
	"fresh-groupbuy/pkg/database"
)

type OrderService struct {
	orderRepo *repository.OrderRepo
}

func NewOrderService() *OrderService {
	return &OrderService{
		orderRepo: repository.NewOrderRepo(),
	}
}

func (s *OrderService) GetByID(ctx context.Context, tx database.DBTX, id, userID int64) (*models.Order, error) {
	if tx == nil {
		tx = database.DB
	}
	order, err := s.orderRepo.GetByID(ctx, tx, id)
	if err != nil {
		return nil, ErrOrderNotFound
	}
	if order.UserID != userID {
		return nil, ErrPermissionDenied
	}
	return order, nil
}

func (s *OrderService) ListByUser(ctx context.Context, tx database.DBTX, userID int64, status string, page, pageSize int) ([]*models.Order, int64, error) {
	if tx == nil {
		tx = database.DB
	}
	offset := (page - 1) * pageSize
	return s.orderRepo.ListByUser(ctx, tx, userID, status, offset, pageSize)
}

func (s *OrderService) Pay(ctx context.Context, tx database.DBTX, id, userID int64) error {
	if tx == nil {
		tx = database.DB
	}
	order, err := s.orderRepo.GetByID(ctx, tx, id)
	if err != nil {
		return ErrOrderNotFound
	}
	if order.UserID != userID {
		return ErrPermissionDenied
	}
	if order.PaymentStatus == models.PaymentStatusPaid {
		return nil
	}
	return s.orderRepo.MarkAsPaid(ctx, tx, id)
}

func (s *OrderService) Ship(ctx context.Context, tx database.DBTX, id int64, trackingNo string) error {
	if tx == nil {
		tx = database.DB
	}
	_, err := s.orderRepo.GetByID(ctx, tx, id)
	if err != nil {
		return ErrOrderNotFound
	}
	_, err = tx.Exec(ctx, `
		UPDATE orders SET status = $1, tracking_no = $2, shipped_at = NOW(), updated_at = NOW()
		WHERE id = $3 AND status = $4
	`, models.OrderStatusShipped, trackingNo, id, models.OrderStatusGrouped)
	return err
}
