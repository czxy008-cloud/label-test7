package repository

import (
	"context"
	"fresh-groupbuy/internal/models"
	"fresh-groupbuy/pkg/database"
)

type ProductRepo struct{}

func NewProductRepo() *ProductRepo {
	return &ProductRepo{}
}

func (r *ProductRepo) GetByID(ctx context.Context, tx database.DBTX, id int64) (*models.Product, error) {
	if tx == nil {
		tx = database.DB
	}
	var p models.Product
	err := tx.QueryRow(ctx, `
		SELECT id, name, description, image_url, original_price, group_price,
		       stock, group_threshold, category, leader_id, status, created_at, updated_at
		FROM products WHERE id = $1 AND status = 'active'
	`, id).Scan(
		&p.ID, &p.Name, &p.Description, &p.ImageURL, &p.OriginalPrice, &p.GroupPrice,
		&p.Stock, &p.GroupThreshold, &p.Category, &p.LeaderID, &p.Status, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepo) List(ctx context.Context, tx database.DBTX, category string, offset, limit int) ([]*models.Product, int64, error) {
	if tx == nil {
		tx = database.DB
	}
	var products []*models.Product
	var total int64

	countSQL := "SELECT COUNT(*) FROM products WHERE status = 'active'"
	listSQL := `
		SELECT id, name, description, image_url, original_price, group_price,
		       stock, group_threshold, category, leader_id, status, created_at, updated_at
		FROM products WHERE status = 'active'
	`
	args := []interface{}{}
	if category != "" {
		countSQL += " AND category = $1"
		listSQL += " AND category = $1"
		args = append(args, category)
	}

	err := tx.QueryRow(ctx, countSQL, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	if category != "" {
		listSQL += " ORDER BY created_at DESC OFFSET $2 LIMIT $3"
	} else {
		listSQL += " ORDER BY created_at DESC OFFSET $1 LIMIT $2"
	}
	args = append(args, offset, limit)

	rows, err := tx.Query(ctx, listSQL, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Product
		err := rows.Scan(
			&p.ID, &p.Name, &p.Description, &p.ImageURL, &p.OriginalPrice, &p.GroupPrice,
			&p.Stock, &p.GroupThreshold, &p.Category, &p.LeaderID, &p.Status, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		products = append(products, &p)
	}

	return products, total, nil
}

func (r *ProductRepo) DecreaseStock(ctx context.Context, tx database.DBTX, id int64, quantity int) error {
	if tx == nil {
		tx = database.DB
	}
	_, err := tx.Exec(ctx, `
		UPDATE products SET stock = stock - $1, updated_at = NOW()
		WHERE id = $2 AND stock >= $1
	`, quantity, id)
	return err
}
