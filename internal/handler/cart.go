package handler

import (
	"fresh-groupbuy/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	cartService *service.CartService
}

func NewCartHandler() *CartHandler {
	return &CartHandler{
		cartService: service.NewCartService(),
	}
}

type addCartRequest struct {
	ProductID int64 `json:"product_id" binding:"required"`
	Quantity  int   `json:"quantity"`
}

type updateCartRequest struct {
	Quantity int `json:"quantity" binding:"required"`
}

func (h *CartHandler) Add(c *gin.Context) {
	var req addCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "参数错误: "+err.Error())
		return
	}

	userID := getUserID(c)

	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	err := h.cartService.Add(c.Request.Context(), nil, userID, req.ProductID, req.Quantity)
	if err != nil {
		handleError(c, err)
		return
	}

	success(c, nil)
}

func (h *CartHandler) Update(c *gin.Context) {
	productIDStr := c.Param("product_id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		fail(c, 400, "无效的商品ID")
		return
	}

	var req updateCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, 400, "参数错误: "+err.Error())
		return
	}

	userID := getUserID(c)

	err = h.cartService.UpdateQuantity(c.Request.Context(), nil, userID, productID, req.Quantity)
	if err != nil {
		handleError(c, err)
		return
	}

	success(c, nil)
}

func (h *CartHandler) Remove(c *gin.Context) {
	productIDStr := c.Param("product_id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		fail(c, 400, "无效的商品ID")
		return
	}

	userID := getUserID(c)

	err = h.cartService.Remove(c.Request.Context(), nil, userID, productID)
	if err != nil {
		handleError(c, err)
		return
	}

	success(c, nil)
}

func (h *CartHandler) List(c *gin.Context) {
	userID := getUserID(c)

	carts, err := h.cartService.List(c.Request.Context(), nil, userID)
	if err != nil {
		handleError(c, err)
		return
	}

	success(c, carts)
}

func (h *CartHandler) Clear(c *gin.Context) {
	userID := getUserID(c)

	err := h.cartService.Clear(c.Request.Context(), nil, userID)
	if err != nil {
		handleError(c, err)
		return
	}

	success(c, nil)
}
