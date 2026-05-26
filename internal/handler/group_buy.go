package handler

import (
	"fresh-groupbuy/internal/service"

	"github.com/gin-gonic/gin"
)

type GroupBuyHandler struct {
	groupBuyService *service.GroupBuyService
}

func NewGroupBuyHandler() *GroupBuyHandler {
	return &GroupBuyHandler{
		groupBuyService: service.NewGroupBuyService(),
	}
}

type createGroupBuyRequest struct {
	ProductID       int64  `json:"product_id" binding:"required"`
	Quantity        int    `json:"quantity"`
	DeliveryAddress string `json:"delivery_address"`
	DeliveryPhone   string `json:"delivery_phone"`
	DeliveryName    string `json:"delivery_name"`
	Remark          string `json:"remark"`
}

type joinGroupBuyRequest struct {
	GroupCode       string `json:"group_code" binding:"required"`
	Quantity        int    `json:"quantity"`
	DeliveryAddress string `json:"delivery_address"`
	DeliveryPhone   string `json:"delivery_phone"`
	DeliveryName    string `json:"delivery_name"`
	Remark          string `json:"remark"`
}

func (h *GroupBuyHandler) Create(c *gin.Context) {
	var req createGroupBuyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "参数错误: "+err.Error())
		return
	}

	userID := getUserID(c)

	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	result, err := h.groupBuyService.CreateGroupBuy(c.Request.Context(), &service.CreateGroupBuyRequest{
		UserID:          userID,
		ProductID:       req.ProductID,
		Quantity:        req.Quantity,
		DeliveryAddress: req.DeliveryAddress,
		DeliveryPhone:   req.DeliveryPhone,
		DeliveryName:    req.DeliveryName,
		Remark:          req.Remark,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	success(c, result)
}

func (h *GroupBuyHandler) Join(c *gin.Context) {
	var req joinGroupBuyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "参数错误: "+err.Error())
		return
	}

	userID := getUserID(c)

	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	result, err := h.groupBuyService.JoinGroupBuy(c.Request.Context(), &service.JoinGroupBuyRequest{
		UserID:          userID,
		GroupCode:       req.GroupCode,
		Quantity:        req.Quantity,
		DeliveryAddress: req.DeliveryAddress,
		DeliveryPhone:   req.DeliveryPhone,
		DeliveryName:    req.DeliveryName,
		Remark:          req.Remark,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	success(c, result)
}

func (h *GroupBuyHandler) GetByCode(c *gin.Context) {
	code := c.Param("code")

	gb, members, err := h.groupBuyService.GetByCode(c.Request.Context(), code)
	if err != nil {
		handleError(c, err)
		return
	}

	success(c, gin.H{
		"group_buy": gb,
		"members":   members,
	})
}

func (h *GroupBuyHandler) List(c *gin.Context) {
	userID := getUserID(c)
	page, pageSize := getPageParams(c)

	groupBuys, total, err := h.groupBuyService.ListByUser(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		handleError(c, err)
		return
	}

	success(c, gin.H{
		"list":      groupBuys,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
