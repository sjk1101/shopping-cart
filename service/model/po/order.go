package po

import (
	"time"

	"github.com/shopspring/decimal"
)

type ShopeeCompletedOrder struct {
	ID               int64           `gorm:"column:id;type:bigint;primary_key;NOT NULL;comment:id"`
	OrderID          string          `gorm:"column:order_id;type:varchar(255);NOT NULL;comment:訂單編號"`
	OrderCreatedAt   time.Time       `gorm:"column:order_created_at;type:DATETIME(6);default:CURRENT_TIMESTAMP(6);comment:訂單成立時間"`
	IsEstablished    bool            `gorm:"column:is_established;type:boolean;NOT NULL;comment:訂單是否成立"`
	OrderCompletedAt *time.Time      `gorm:"column:order_completed_at;type:DATETIME(6);NULL;comment:訂單完成時間"`
	AllocateAt       *time.Time      `gorm:"column:allocate_at;type:DATETIME(6);NULL;comment:撥款日"`
	Price            decimal.Decimal `gorm:"column:price;type:decimal(20,4);NOT NULL;comment:商品金額"`
	CouponDiscount   decimal.Decimal `gorm:"column:coupon_discount;type:decimal(10,4);NOT NULL;comment:賣場優惠券"`
	DealFee          decimal.Decimal `gorm:"column:deal_fee;type:decimal(10,4);NOT NULL;comment:成交手續費"`
	ActivityFee      decimal.Decimal `gorm:"column:activity_fee;type:decimal(10,4);NOT NULL;comment:活動服務費"`
	CashFlowCost     decimal.Decimal `gorm:"column:cash_flow_cost;type:decimal(10,4);NOT NULL;comment:金流服務費"`
	CreatedAt        *time.Time      `gorm:"column:created_at;type:DATETIME(6);default:CURRENT_TIMESTAMP(6);comment:建立時間"`
	UpdatedAt        *time.Time      `gorm:"column:updated_at;type:DATETIME(6);default:CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6);comment:更新時間"`
}

func (o *ShopeeCompletedOrder) TableName() string {
	return "shopee_completed_order"
}
