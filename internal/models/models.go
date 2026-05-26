package models

import "time"

type User struct {
	ID           int64     `json:"id"`
	Phone        string    `json:"phone"`
	Nickname     string    `json:"nickname"`
	AvatarURL    string    `json:"avatar_url"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Leader struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	CommunityName string    `json:"community_name"`
	Address       string    `json:"address"`
	ContactPhone  string    `json:"contact_phone"`
	Status        string    `json:"status"`
	AuditNote     string    `json:"audit_note"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Product struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	ImageURL        string    `json:"image_url"`
	OriginalPrice   float64   `json:"original_price"`
	GroupPrice      float64   `json:"group_price"`
	Stock           int       `json:"stock"`
	GroupThreshold  int       `json:"group_threshold"`
	Category        string    `json:"category"`
	LeaderID        int64     `json:"leader_id"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Cart struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	ProductID int64     `json:"product_id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GroupBuy struct {
	ID           int64     `json:"id"`
	GroupCode    string    `json:"group_code"`
	ProductID    int64     `json:"product_id"`
	InitiatorID  int64     `json:"initiator_id"`
	LeaderID     int64     `json:"leader_id"`
	CurrentCount int       `json:"current_count"`
	TargetCount  int       `json:"target_count"`
	Status       string    `json:"status"`
	ExpireAt     time.Time `json:"expire_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type GroupBuyMember struct {
	ID          int64     `json:"id"`
	GroupBuyID  int64     `json:"group_buy_id"`
	UserID      int64     `json:"user_id"`
	OrderID     int64     `json:"order_id"`
	IsInitiator bool      `json:"is_initiator"`
	JoinedAt    time.Time `json:"joined_at"`
}

type Order struct {
	ID              int64      `json:"id"`
	OrderNo         string     `json:"order_no"`
	UserID          int64      `json:"user_id"`
	ProductID       int64      `json:"product_id"`
	GroupBuyID      int64      `json:"group_buy_id"`
	LeaderID        int64      `json:"leader_id"`
	UnitPrice       float64    `json:"unit_price"`
	Quantity        int        `json:"quantity"`
	TotalAmount     float64    `json:"total_amount"`
	Status          string     `json:"status"`
	PaymentStatus   string     `json:"payment_status"`
	DeliveryAddress string     `json:"delivery_address"`
	DeliveryPhone   string     `json:"delivery_phone"`
	DeliveryName    string     `json:"delivery_name"`
	TrackingNo      string     `json:"tracking_no"`
	Remark          string     `json:"remark"`
	PaidAt          *time.Time `json:"paid_at"`
	ShippedAt       *time.Time `json:"shipped_at"`
	DeliveredAt     *time.Time `json:"delivered_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

const (
	OrderStatusPendingGroup = "pending_group"
	OrderStatusGrouped      = "grouped"
	OrderStatusShipped      = "shipped"
	OrderStatusDelivered    = "delivered"
	OrderStatusCancelled    = "cancelled"

	GroupBuyStatusPending  = "pending"
	GroupBuyStatusSuccess  = "success"
	GroupBuyStatusFailed   = "failed"
	GroupBuyStatusCancelled = "cancelled"

	PaymentStatusUnpaid   = "unpaid"
	PaymentStatusPaid     = "paid"
	PaymentStatusRefunded = "refunded"
)
