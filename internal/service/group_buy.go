package service

import (
	"context"
	"fmt"
	"fresh-groupbuy/internal/models"
	"fresh-groupbuy/internal/repository"
	"fresh-groupbuy/pkg/database"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type GroupBuyService struct {
	groupBuyRepo *repository.GroupBuyRepo
	productRepo  *repository.ProductRepo
	orderRepo    *repository.OrderRepo
}

func NewGroupBuyService() *GroupBuyService {
	return &GroupBuyService{
		groupBuyRepo: repository.NewGroupBuyRepo(),
		productRepo:  repository.NewProductRepo(),
		orderRepo:    repository.NewOrderRepo(),
	}
}

type CreateGroupBuyRequest struct {
	UserID          int64
	ProductID       int64
	Quantity        int
	DeliveryAddress string
	DeliveryPhone   string
	DeliveryName    string
	Remark          string
}

type GroupBuyResult struct {
	GroupBuyID   int64  `json:"group_buy_id"`
	GroupCode    string `json:"group_code"`
	OrderID      int64  `json:"order_id"`
	OrderNo      string `json:"order_no"`
	CurrentCount int    `json:"current_count"`
	TargetCount  int    `json:"target_count"`
	Status       string `json:"status"`
}

func (s *GroupBuyService) CreateGroupBuy(ctx context.Context, req *CreateGroupBuyRequest) (*GroupBuyResult, error) {
	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	product, err := s.productRepo.GetByID(ctx, nil, req.ProductID)
	if err != nil {
		return nil, ErrProductNotFound
	}

	if product.Stock < req.Quantity {
		return nil, ErrInsufficientStock
	}

	expireHours := 24
	if hoursStr := os.Getenv("GROUP_BUY_EXPIRE_HOURS"); hoursStr != "" {
		if h, err := strconv.Atoi(hoursStr); err == nil {
			expireHours = h
		}
	}

	groupCode := generateGroupCode()

	gb := &models.GroupBuy{
		GroupCode:   groupCode,
		ProductID:   product.ID,
		InitiatorID: req.UserID,
		LeaderID:    product.LeaderID,
		TargetCount: product.GroupThreshold,
		ExpireAt:    time.Now().Add(time.Duration(expireHours) * time.Hour),
	}

	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	groupBuyID, err := s.groupBuyRepo.Create(ctx, tx, gb)
	if err != nil {
		return nil, fmt.Errorf("创建拼团失败: %w", err)
	}

	orderNo := generateOrderNo()
	totalAmount := float64(req.Quantity) * product.GroupPrice
	order := &models.Order{
		OrderNo:         orderNo,
		UserID:          req.UserID,
		ProductID:       product.ID,
		GroupBuyID:      groupBuyID,
		LeaderID:        product.LeaderID,
		UnitPrice:       product.GroupPrice,
		Quantity:        req.Quantity,
		TotalAmount:     totalAmount,
		Status:          models.OrderStatusPendingGroup,
		PaymentStatus:   models.PaymentStatusUnpaid,
		DeliveryAddress: req.DeliveryAddress,
		DeliveryPhone:   req.DeliveryPhone,
		DeliveryName:    req.DeliveryName,
		Remark:          req.Remark,
	}

	orderID, err := s.orderRepo.Create(ctx, tx, order)
	if err != nil {
		return nil, fmt.Errorf("创建订单失败: %w", err)
	}

	err = s.groupBuyRepo.AddMember(ctx, tx, groupBuyID, req.UserID, orderID, true)
	if err != nil {
		return nil, fmt.Errorf("添加拼团成员失败: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &GroupBuyResult{
		GroupBuyID:   groupBuyID,
		GroupCode:    groupCode,
		OrderID:      orderID,
		OrderNo:      orderNo,
		CurrentCount: 1,
		TargetCount:  product.GroupThreshold,
		Status:       models.GroupBuyStatusPending,
	}, nil
}

type JoinGroupBuyRequest struct {
	UserID          int64
	GroupCode       string
	Quantity        int
	DeliveryAddress string
	DeliveryPhone   string
	DeliveryName    string
	Remark          string
}

func (s *GroupBuyService) JoinGroupBuy(ctx context.Context, req *JoinGroupBuyRequest) (*GroupBuyResult, error) {
	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	gb, err := s.groupBuyRepo.GetByCode(ctx, nil, req.GroupCode)
	if err != nil {
		return nil, ErrGroupBuyNotFound
	}

	if gb.Status != models.GroupBuyStatusPending {
		return nil, ErrGroupBuyClosed
	}

	if gb.ExpireAt.Before(time.Now()) {
		return nil, ErrGroupBuyClosed
	}

	hasMember, err := s.groupBuyRepo.HasMember(ctx, nil, gb.ID, req.UserID)
	if err != nil {
		return nil, err
	}
	if hasMember {
		return nil, ErrAlreadyInGroupBuy
	}

	if gb.CurrentCount >= gb.TargetCount {
		return nil, ErrGroupBuyFull
	}

	product, err := s.productRepo.GetByID(ctx, nil, gb.ProductID)
	if err != nil {
		return nil, ErrProductNotFound
	}

	if product.Stock < req.Quantity {
		return nil, ErrInsufficientStock
	}

	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	orderNo := generateOrderNo()
	totalAmount := float64(req.Quantity) * product.GroupPrice
	order := &models.Order{
		OrderNo:         orderNo,
		UserID:          req.UserID,
		ProductID:       product.ID,
		GroupBuyID:      gb.ID,
		LeaderID:        gb.LeaderID,
		UnitPrice:       product.GroupPrice,
		Quantity:        req.Quantity,
		TotalAmount:     totalAmount,
		Status:          models.OrderStatusPendingGroup,
		PaymentStatus:   models.PaymentStatusUnpaid,
		DeliveryAddress: req.DeliveryAddress,
		DeliveryPhone:   req.DeliveryPhone,
		DeliveryName:    req.DeliveryName,
		Remark:          req.Remark,
	}

	orderID, err := s.orderRepo.Create(ctx, tx, order)
	if err != nil {
		return nil, fmt.Errorf("创建订单失败: %w", err)
	}

	err = s.groupBuyRepo.AddMember(ctx, tx, gb.ID, req.UserID, orderID, false)
	if err != nil {
		return nil, fmt.Errorf("添加拼团成员失败: %w", err)
	}

	newCount, err := s.groupBuyRepo.IncrementCount(ctx, tx, gb.ID)
	if err != nil {
		return nil, fmt.Errorf("更新拼团人数失败: %w", err)
	}

	if newCount >= gb.TargetCount {
		err = s.processGroupSuccess(ctx, tx, gb, product)
		if err != nil {
			return nil, fmt.Errorf("处理成团失败: %w", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	status := models.GroupBuyStatusPending
	if newCount >= gb.TargetCount {
		status = models.GroupBuyStatusSuccess
	}

	return &GroupBuyResult{
		GroupBuyID:   gb.ID,
		GroupCode:    gb.GroupCode,
		OrderID:      orderID,
		OrderNo:      orderNo,
		CurrentCount: newCount,
		TargetCount:  gb.TargetCount,
		Status:       status,
	}, nil
}

func (s *GroupBuyService) processGroupSuccess(ctx context.Context, tx pgx.Tx, gb *models.GroupBuy, product *models.Product) error {
	err := s.groupBuyRepo.UpdateStatus(ctx, tx, gb.ID, models.GroupBuyStatusSuccess)
	if err != nil {
		return err
	}

	members, err := s.groupBuyRepo.GetMembers(ctx, tx, gb.ID)
	if err != nil {
		return err
	}

	totalQuantity := 0
	for _, member := range members {
		order, err := s.orderRepo.GetByID(ctx, tx, member.OrderID)
		if err != nil {
			return err
		}
		totalQuantity += order.Quantity
	}

	err = s.productRepo.DecreaseStock(ctx, tx, product.ID, totalQuantity)
	if err != nil {
		return err
	}

	err = s.orderRepo.UpdateStatusByGroupBuyID(ctx, tx, gb.ID, models.OrderStatusGrouped)
	if err != nil {
		return err
	}

	return nil
}

func (s *GroupBuyService) GetByCode(ctx context.Context, code string) (*models.GroupBuy, []*models.GroupBuyMember, error) {
	gb, err := s.groupBuyRepo.GetByCode(ctx, nil, code)
	if err != nil {
		return nil, nil, ErrGroupBuyNotFound
	}

	members, err := s.groupBuyRepo.GetMembers(ctx, nil, gb.ID)
	if err != nil {
		return nil, nil, err
	}

	return gb, members, nil
}

func (s *GroupBuyService) ListByUser(ctx context.Context, userID int64, page, pageSize int) ([]*models.GroupBuy, int64, error) {
	offset := (page - 1) * pageSize
	return s.groupBuyRepo.ListByUser(ctx, nil, userID, offset, pageSize)
}

func generateGroupCode() string {
	return uuid.New().String()[:8]
}

func generateOrderNo() string {
	now := time.Now()
	return fmt.Sprintf("GB%s%06d", now.Format("20060102150405"), time.Now().UnixNano()%1000000)
}
