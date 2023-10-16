package bo

import (
	"time"

	"github.com/shopspring/decimal"
)

type ShopeeOrder struct {
	OrderID          string
	OrderCreatedAt   time.Time
	IsEstablished    bool
	OrderCompletedAt *time.Time
	AllocateAt       *time.Time
	TotalPrice       decimal.Decimal
	CouponDiscount   decimal.Decimal
	DealFee          decimal.Decimal
	ActivityFee      decimal.Decimal
	CashFlowCost     decimal.Decimal
}

type ShopeeOrderDetail struct {
	OrderID          string
	OrderCreatedAt   time.Time
	IsEstablished    bool
	OrderCompletedAt *time.Time
	Product          string
	Price            decimal.Decimal
	Quantity         int
	CouponDiscount   decimal.Decimal
	DealFee          decimal.Decimal
	ActivityFee      decimal.Decimal
	CashFlowCost     decimal.Decimal
}
