package repository

import (
	"context"
	"fresh-groupbuy/internal/models"
	"fresh-groupbuy/pkg/database"
)

type GroupBuyRepo struct{}

func NewGroupBuyRepo() *GroupBuyRepo {
	return &GroupBuyRepo{}
}

func (r *GroupBuyRepo) Create(ctx context.Context, tx database.DBTX, gb *models.GroupBuy) (int64, error) {
	if tx == nil {
		tx = database.DB
	}
	var id int64
	err := tx.QueryRow(ctx, `
		INSERT INTO group_buys (group_code, product_id, initiator_id, leader_id, current_count, target_count, status, expire_at)
		VALUES ($1, $2, $3, $4, 1, $5, 'pending', $6)
		RETURNING id
	`, gb.GroupCode, gb.ProductID, gb.InitiatorID, gb.LeaderID, gb.TargetCount, gb.ExpireAt).Scan(&id)
	return id, err
}

func (r *GroupBuyRepo) GetByCode(ctx context.Context, tx database.DBTX, code string) (*models.GroupBuy, error) {
	if tx == nil {
		tx = database.DB
	}
	var gb models.GroupBuy
	err := tx.QueryRow(ctx, `
		SELECT id, group_code, product_id, initiator_id, leader_id, current_count,
		       target_count, status, expire_at, created_at, updated_at
		FROM group_buys WHERE group_code = $1
	`, code).Scan(
		&gb.ID, &gb.GroupCode, &gb.ProductID, &gb.InitiatorID, &gb.LeaderID,
		&gb.CurrentCount, &gb.TargetCount, &gb.Status, &gb.ExpireAt, &gb.CreatedAt, &gb.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &gb, nil
}

func (r *GroupBuyRepo) GetByID(ctx context.Context, tx database.DBTX, id int64) (*models.GroupBuy, error) {
	if tx == nil {
		tx = database.DB
	}
	var gb models.GroupBuy
	err := tx.QueryRow(ctx, `
		SELECT id, group_code, product_id, initiator_id, leader_id, current_count,
		       target_count, status, expire_at, created_at, updated_at
		FROM group_buys WHERE id = $1
	`, id).Scan(
		&gb.ID, &gb.GroupCode, &gb.ProductID, &gb.InitiatorID, &gb.LeaderID,
		&gb.CurrentCount, &gb.TargetCount, &gb.Status, &gb.ExpireAt, &gb.CreatedAt, &gb.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &gb, nil
}

func (r *GroupBuyRepo) AddMember(ctx context.Context, tx database.DBTX, groupBuyID, userID, orderID int64, isInitiator bool) error {
	if tx == nil {
		tx = database.DB
	}
	_, err := tx.Exec(ctx, `
		INSERT INTO group_buy_members (group_buy_id, user_id, order_id, is_initiator, joined_at)
		VALUES ($1, $2, $3, $4, NOW())
	`, groupBuyID, userID, orderID, isInitiator)
	return err
}

func (r *GroupBuyRepo) IncrementCount(ctx context.Context, tx database.DBTX, groupBuyID int64) (int, error) {
	if tx == nil {
		tx = database.DB
	}
	var currentCount int
	err := tx.QueryRow(ctx, `
		UPDATE group_buys SET current_count = current_count + 1, updated_at = NOW()
		WHERE id = $1
		RETURNING current_count
	`, groupBuyID).Scan(&currentCount)
	return currentCount, err
}

func (r *GroupBuyRepo) UpdateStatus(ctx context.Context, tx database.DBTX, groupBuyID int64, status string) error {
	if tx == nil {
		tx = database.DB
	}
	_, err := tx.Exec(ctx, `
		UPDATE group_buys SET status = $1, updated_at = NOW() WHERE id = $2
	`, status, groupBuyID)
	return err
}

func (r *GroupBuyRepo) HasMember(ctx context.Context, tx database.DBTX, groupBuyID, userID int64) (bool, error) {
	if tx == nil {
		tx = database.DB
	}
	var exists bool
	err := tx.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM group_buy_members WHERE group_buy_id = $1 AND user_id = $2)
	`, groupBuyID, userID).Scan(&exists)
	return exists, err
}

func (r *GroupBuyRepo) GetMembers(ctx context.Context, tx database.DBTX, groupBuyID int64) ([]*models.GroupBuyMember, error) {
	if tx == nil {
		tx = database.DB
	}
	rows, err := tx.Query(ctx, `
		SELECT id, group_buy_id, user_id, order_id, is_initiator, joined_at
		FROM group_buy_members WHERE group_buy_id = $1 ORDER BY joined_at ASC
	`, groupBuyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*models.GroupBuyMember
	for rows.Next() {
		var m models.GroupBuyMember
		err := rows.Scan(&m.ID, &m.GroupBuyID, &m.UserID, &m.OrderID, &m.IsInitiator, &m.JoinedAt)
		if err != nil {
			return nil, err
		}
		members = append(members, &m)
	}
	return members, nil
}

func (r *GroupBuyRepo) ListByUser(ctx context.Context, tx database.DBTX, userID int64, offset, limit int) ([]*models.GroupBuy, int64, error) {
	if tx == nil {
		tx = database.DB
	}
	var groupBuys []*models.GroupBuy
	var total int64

	err := tx.QueryRow(ctx, `
		SELECT COUNT(DISTINCT gb.id) FROM group_buys gb
		INNER JOIN group_buy_members gbm ON gb.id = gbm.group_buy_id
		WHERE gbm.user_id = $1
	`, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := tx.Query(ctx, `
		SELECT DISTINCT gb.id, gb.group_code, gb.product_id, gb.initiator_id, gb.leader_id,
		       gb.current_count, gb.target_count, gb.status, gb.expire_at, gb.created_at, gb.updated_at
		FROM group_buys gb
		INNER JOIN group_buy_members gbm ON gb.id = gbm.group_buy_id
		WHERE gbm.user_id = $1
		ORDER BY gb.created_at DESC
		OFFSET $2 LIMIT $3
	`, userID, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var gb models.GroupBuy
		err := rows.Scan(
			&gb.ID, &gb.GroupCode, &gb.ProductID, &gb.InitiatorID, &gb.LeaderID,
			&gb.CurrentCount, &gb.TargetCount, &gb.Status, &gb.ExpireAt, &gb.CreatedAt, &gb.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		groupBuys = append(groupBuys, &gb)
	}

	return groupBuys, total, nil
}
