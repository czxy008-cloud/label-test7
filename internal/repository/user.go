package repository

import (
	"context"
	"fresh-groupbuy/internal/models"
	"fresh-groupbuy/pkg/database"
)

type UserRepo struct{}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (r *UserRepo) GetByID(ctx context.Context, tx database.DBTX, id int64) (*models.User, error) {
	if tx == nil {
		tx = database.DB
	}
	var u models.User
	err := tx.QueryRow(ctx, `
		SELECT id, phone, nickname, avatar_url, password_hash, role, created_at, updated_at
		FROM users WHERE id = $1
	`, id).Scan(
		&u.ID, &u.Phone, &u.Nickname, &u.AvatarURL, &u.PasswordHash,
		&u.Role, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) GetByPhone(ctx context.Context, tx database.DBTX, phone string) (*models.User, error) {
	if tx == nil {
		tx = database.DB
	}
	var u models.User
	err := tx.QueryRow(ctx, `
		SELECT id, phone, nickname, avatar_url, password_hash, role, created_at, updated_at
		FROM users WHERE phone = $1
	`, phone).Scan(
		&u.ID, &u.Phone, &u.Nickname, &u.AvatarURL, &u.PasswordHash,
		&u.Role, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
