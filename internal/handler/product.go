package handler

import (
	"fresh-groupbuy/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{
		productService: service.NewProductService(),
	}
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		fail(c, 400, "无效的商品ID")
		return
	}

	product, err := h.productService.GetByID(c.Request.Context(), nil, id)
	if err != nil {
		handleError(c, err)
		return
	}

	success(c, product)
}

func (h *ProductHandler) List(c *gin.Context) {
	category := c.Query("category")
	page, pageSize := getPageParams(c)

	products, total, err := h.productService.List(c.Request.Context(), nil, category, page, pageSize)
	if err != nil {
		handleError(c, err)
		return
	}

	success(c, gin.H{
		"list":      products,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
