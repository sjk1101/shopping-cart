package po

import (
	"time"

	"github.com/shopspring/decimal"
)

type Product struct {
	Name      string          `gorm:"column:name;type:varchar(512);NOT NULL;comment:產品名稱"`
	Number    string          `gorm:"column:number;type:varchar(16);Default:'';comment:產品貨號"`
	Cost      decimal.Decimal `gorm:"column:cost;type:decimal(10,4);comment:產品成本"`
	Quantity  int             `gorm:"column:quantity;type:int;comment:庫存"`
	Image     string          `gorm:"column:image;type:varchar(255);comment:圖片路徑"`
	CreatedAt *time.Time      `gorm:"column:created_at;type:DATETIME(6);default:CURRENT_TIMESTAMP(6);comment:建立時間"`
	UpdatedAt *time.Time      `gorm:"column:updated_at;type:DATETIME(6);default:CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6);comment:更新時間"`
}

func (p *Product) TableName() string {
	return "product"
}
