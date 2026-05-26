package repository

import (
	"context"
	"fresh-groupbuy/internal/models"
	"fresh-groupbuy/pkg/database"
)

type CartRepo struct{}

func NewCartRepo() *CartRepo {
	return &CartRepo{}
}

func (r *CartRepo) Add(ctx context.Context, tx database.DBTX, userID, productID int64, quantity int) error {
	if tx == nil {
		tx = database.DB
	}
	_, err := tx.Exec(ctx, `
		INSERT INTO carts (user_id, product_id, quantity, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		ON CONFLICT (user_id, product_id) DO UPDATE
		SET quantity = carts.quantity + $3, updated_at = NOW()
	`, userID, productID, quantity)
	return err
}

func (r *CartRepo) UpdateQuantity(ctx context.Context, tx database.DBTX, userID, productID int64, quantity int) error {
	if tx == nil {
		tx = database.DB
	}
	_, err := tx.Exec(ctx, `
		UPDATE carts SET quantity = $1, updated_at = NOW()
		WHERE user_id = $2 AND product_id = $3
	`, quantity, userID, productID)
	return err
}

func (r *CartRepo) Remove(ctx context.Context, tx database.DBTX, userID, productID int64) error {
	if tx == nil {
		tx = database.DB
	}
	_, err := tx.Exec(ctx, `
		DELETE FROM carts WHERE user_id = $1 AND product_id = $2
	`, userID, productID)
	return err
}

func (r *CartRepo) List(ctx context.Context, tx database.DBTX, userID int64) ([]*models.Cart, error) {
	if tx == nil {
		tx = database.DB
	}
	rows, err := tx.Query(ctx, `
		SELECT id, user_id, product_id, quantity, created_at, updated_at
		FROM carts WHERE user_id = $1 ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var carts []*models.Cart
	for rows.Next() {
		var c models.Cart
		err := rows.Scan(&c.ID, &c.UserID, &c.ProductID, &c.Quantity, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		carts = append(carts, &c)
	}
	return carts, nil
}

func (r *CartRepo) Clear(ctx context.Context, tx database.DBTX, userID int64) error {
	if tx == nil {
		tx = database.DB
	}
	_, err := tx.Exec(ctx, "DELETE FROM carts WHERE user_id = $1", userID)
	return err
}
