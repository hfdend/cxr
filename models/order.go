package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// 订单状态
// 1: 等待付款
// 2: 付款成功
// 3: 订单发货中 (月够订单)
// 4: 发货完成
// swagger:model OrderStatus
type OrderStatus int

const (
	OrderStatusWatting OrderStatus = iota + 1
	OrderStatusSuccess
	OrderStatusDelivering
	OrderStatusDeliveried
)

// 订单
// swagger:model Order
type Order struct {
	Model
	OrderID string `json:"order_id"`
	UserId  int    `json:"user_id"`
	// 如果是月够的则有值
	PlanId int `json:"plan_id"`
	// 金额
	Price float64 `json:"price"`
	// 手续费
	Fee float64 `json:"fee"`
	// 支付金额
	PaymentPrice float64 `json:"payment_price"`
	// 支付方式 1: 微信支付
	PaymentMethod int `json:"payment_method"`
	// 买家留言
	Notice string      `json:"notice"`
	Status OrderStatus `json:"status"`
	// 创建时间
	Created int64 `json:"created"`
	// 支付截止时间
	ExpTime int64 `json:"exp_time"`
	// 支付时间
	PaymentTime int64 `json:"payment_time"`
	UpdateTime  int64 `json:"update_time"`

	OrderProducts []*OrderProduct `json:"order_products" gorm:"-"`
	OrderAddress  *OrderAddress   `json:"order_address" gorm:"-"`
}

var OrderDefault Order

func (Order) TableName() string {
	return "order"
}

func (o *Order) Insert(db *gorm.DB) error {
	o.Created = time.Now().Unix()
	o.UpdateTime = time.Now().Unix()
	return db.Create(o).Error
}
