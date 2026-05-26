package handler

import (
	"fresh-groupbuy/internal/models"
	"fresh-groupbuy/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, models.Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

func fail(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, models.Response{
		Code:    code,
		Message: message,
	})
}

func handleError(c *gin.Context, err error) {
	switch err {
	case service.ErrProductNotFound:
		fail(c, 404, "商品不存在")
	case service.ErrInsufficientStock:
		fail(c, 400, "库存不足")
	case service.ErrGroupBuyNotFound:
		fail(c, 404, "拼团不存在")
	case service.ErrGroupBuyClosed:
		fail(c, 400, "拼团已结束")
	case service.ErrAlreadyInGroupBuy:
		fail(c, 400, "已加入该拼团")
	case service.ErrGroupBuyFull:
		fail(c, 400, "拼团人数已满")
	case service.ErrInvalidQuantity:
		fail(c, 400, "无效的数量")
	case service.ErrOrderNotFound:
		fail(c, 404, "订单不存在")
	case service.ErrPermissionDenied:
		fail(c, 403, "无权限操作")
	case service.ErrUserNotFound:
		fail(c, 404, "用户不存在")
	default:
		fail(c, 500, err.Error())
	}
}

func getUserID(c *gin.Context) int64 {
	userID, exists := c.Get("user_id")
	if !exists {
		return 2 // 默认测试用户ID
	}
	return userID.(int64)
}

func getPageParams(c *gin.Context) (int, int) {
	page := 1
	pageSize := 10

	if p := c.Query("page"); p != "" {
		if v, err := parseInt(p); err == nil && v > 0 {
			page = v
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if v, err := parseInt(ps); err == nil && v > 0 && v <= 100 {
			pageSize = v
		}
	}

	return page, pageSize
}

func parseInt(s string) (int, error) {
	var v int
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, nil
		}
		v = v*10 + int(c-'0')
	}
	return v, nil
}
