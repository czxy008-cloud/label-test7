package service

import "errors"

var (
	ErrProductNotFound       = errors.New("商品不存在")
	ErrInsufficientStock     = errors.New("库存不足")
	ErrGroupBuyNotFound      = errors.New("拼团不存在")
	ErrGroupBuyClosed        = errors.New("拼团已结束")
	ErrAlreadyInGroupBuy     = errors.New("已加入该拼团")
	ErrGroupBuyFull          = errors.New("拼团人数已满")
	ErrInvalidQuantity       = errors.New("无效的数量")
	ErrOrderNotFound         = errors.New("订单不存在")
	ErrPermissionDenied      = errors.New("无权限操作")
	ErrUserNotFound          = errors.New("用户不存在")
)
