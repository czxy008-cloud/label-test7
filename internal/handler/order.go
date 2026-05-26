package handler

import (
	"fresh-groupbuy/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{
		orderService: service.NewOrderService(),
	}
}

type shipRequest struct {
	TrackingNo string `json:"tracking_no"`
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		fail(c, 400, "无效的订单ID")
		return
	}

	userID := getUserID(c)

	order, err := h.orderService.GetByID(c.Request.Context(), nil, id, userID)
	if err != nil {
		handleError(c, err)
		return
	}

	success(c, order)
}

func (h *OrderHandler) List(c *gin.Context) {
	status := c.Query("status")
	page, pageSize := getPageParams(c)
	userID := getUserID(c)

	orders, total, err := h.orderService.ListByUser(c.Request.Context(), nil, userID, status, page, pageSize)
	if err != nil {
		handleError(c, err)
		return
	}

	success(c, gin.H{
		"list":      orders,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *OrderHandler) Pay(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		fail(c, 400, "无效的订单ID")
		return
	}

	userID := getUserID(c)

	err = h.orderService.Pay(c.Request.Context(), nil, id, userID)
	if err != nil {
		handleError(c, err)
		return
	}

	success(c, nil)
}

func (h *OrderHandler) Ship(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		fail(c, 400, "无效的订单ID")
		return
	}

	var req shipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "参数错误: "+err.Error())
		return
	}

	err = h.orderService.Ship(c.Request.Context(), nil, id, req.TrackingNo)
	if err != nil {
		handleError(c, err)
		return
	}

	success(c, nil)
}
