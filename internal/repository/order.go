package repository

import (
	"context"
	"fresh-groupbuy/internal/models"
	"fresh-groupbuy/pkg/database"
)

type OrderRepo struct{}

func NewOrderRepo() *OrderRepo {
	return &OrderRepo{}
}

func (r *OrderRepo) Create(ctx context.Context, tx database.DBTX, order *models.Order) (int64, error) {
	if tx == nil {
		tx = database.DB
	}
	var id int64
	err := tx.QueryRow(ctx, `
		INSERT INTO orders (order_no, user_id, product_id, group_buy_id, leader_id,
		                   unit_price, quantity, total_amount, status, payment_status,
		                   delivery_address, delivery_phone, delivery_name, remark)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id
	`, order.OrderNo, order.UserID, order.ProductID, order.GroupBuyID, order.LeaderID,
		order.UnitPrice, order.Quantity, order.TotalAmount, order.Status, order.PaymentStatus,
		order.DeliveryAddress, order.DeliveryPhone, order.DeliveryName, order.Remark,
	).Scan(&id)
	return id, err
}

func (r *OrderRepo) GetByID(ctx context.Context, tx database.DBTX, id int64) (*models.Order, error) {
	if tx == nil {
		tx = database.DB
	}
	var o models.Order
	err := tx.QueryRow(ctx, `
		SELECT id, order_no, user_id, product_id, group_buy_id, leader_id,
		       unit_price, quantity, total_amount, status, payment_status,
		       delivery_address, delivery_phone, delivery_name, tracking_no, remark,
		       paid_at, shipped_at, delivered_at, created_at, updated_at
		FROM orders WHERE id = $1
	`, id).Scan(
		&o.ID, &o.OrderNo, &o.UserID, &o.ProductID, &o.GroupBuyID, &o.LeaderID,
		&o.UnitPrice, &o.Quantity, &o.TotalAmount, &o.Status, &o.PaymentStatus,
		&o.DeliveryAddress, &o.DeliveryPhone, &o.DeliveryName, &o.TrackingNo, &o.Remark,
		&o.PaidAt, &o.ShippedAt, &o.DeliveredAt, &o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *OrderRepo) UpdateStatus(ctx context.Context, tx database.DBTX, id int64, status string) error {
	if tx == nil {
		tx = database.DB
	}
	_, err := tx.Exec(ctx, `
		UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2
	`, status, id)
	return err
}

func (r *OrderRepo) UpdateStatusByGroupBuyID(ctx context.Context, tx database.DBTX, groupBuyID int64, status string) error {
	if tx == nil {
		tx = database.DB
	}
	_, err := tx.Exec(ctx, `
		UPDATE orders SET status = $1, updated_at = NOW() WHERE group_buy_id = $2
	`, status, groupBuyID)
	return err
}

func (r *OrderRepo) MarkAsPaid(ctx context.Context, tx database.DBTX, id int64) error {
	if tx == nil {
		tx = database.DB
	}
	_, err := tx.Exec(ctx, `
		UPDATE orders SET payment_status = 'paid', paid_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`, id)
	return err
}

func (r *OrderRepo) ListByUser(ctx context.Context, tx database.DBTX, userID int64, status string, offset, limit int) ([]*models.Order, int64, error) {
	if tx == nil {
		tx = database.DB
	}
	var orders []*models.Order
	var total int64

	countSQL := "SELECT COUNT(*) FROM orders WHERE user_id = $1"
	listSQL := `
		SELECT id, order_no, user_id, product_id, group_buy_id, leader_id,
		       unit_price, quantity, total_amount, status, payment_status,
		       delivery_address, delivery_phone, delivery_name, tracking_no, remark,
		       paid_at, shipped_at, delivered_at, created_at, updated_at
		FROM orders WHERE user_id = $1
	`

	args := []interface{}{userID}
	argIndex := 2

	if status != "" {
		countSQL += " AND status = $" + string(rune('0'+argIndex))
		listSQL += " AND status = $" + string(rune('0'+argIndex))
		args = append(args, status)
		argIndex++
	}

	err := tx.QueryRow(ctx, countSQL, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	listSQL += " ORDER BY created_at DESC OFFSET $" + string(rune('0'+argIndex)) + " LIMIT $" + string(rune('0'+argIndex+1))
	args = append(args, offset, limit)

	rows, err := tx.Query(ctx, listSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var o models.Order
		err := rows.Scan(
			&o.ID, &o.OrderNo, &o.UserID, &o.ProductID, &o.GroupBuyID, &o.LeaderID,
			&o.UnitPrice, &o.Quantity, &o.TotalAmount, &o.Status, &o.PaymentStatus,
			&o.DeliveryAddress, &o.DeliveryPhone, &o.DeliveryName, &o.TrackingNo, &o.Remark,
			&o.PaidAt, &o.ShippedAt, &o.DeliveredAt, &o.CreatedAt, &o.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		orders = append(orders, &o)
	}

	return orders, total, nil
}

func (r *OrderRepo) GetByGroupBuyAndUser(ctx context.Context, tx database.DBTX, groupBuyID, userID int64) (*models.Order, error) {
	if tx == nil {
		tx = database.DB
	}
	var o models.Order
	err := tx.QueryRow(ctx, `
		SELECT id, order_no, user_id, product_id, group_buy_id, leader_id,
		       unit_price, quantity, total_amount, status, payment_status,
		       delivery_address, delivery_phone, delivery_name, tracking_no, remark,
		       paid_at, shipped_at, delivered_at, created_at, updated_at
		FROM orders WHERE group_buy_id = $1 AND user_id = $2
	`, groupBuyID, userID).Scan(
		&o.ID, &o.OrderNo, &o.UserID, &o.ProductID, &o.GroupBuyID, &o.LeaderID,
		&o.UnitPrice, &o.Quantity, &o.TotalAmount, &o.Status, &o.PaymentStatus,
		&o.DeliveryAddress, &o.DeliveryPhone, &o.DeliveryName, &o.TrackingNo, &o.Remark,
		&o.PaidAt, &o.ShippedAt, &o.DeliveredAt, &o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &o, nil
}
