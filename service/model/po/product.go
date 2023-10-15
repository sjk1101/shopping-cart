package po

import "time"

type Product struct {
	ID        int64      `gorm:"column:id;type:bigint;primary_key;NOT NULL;comment:id"`
	Name      string     `gorm:"column:name;type:varchar(255);NOT NULL;comment:產品名稱"`
	Amount    int        `gorm:"column:amount;type:int;comment:產品價格"`
	Inventory int        `gorm:"column:inventory;type:int;comment:庫存"`
	Image     string     `gorm:"column:image;type:varchar(255);comment:圖片路徑"`
	CreatedAt *time.Time `gorm:"column:created_at;type:DATETIME(6);default:CURRENT_TIMESTAMP(6);comment:建立時間"`
	UpdatedAt *time.Time `gorm:"column:updated_at;type:DATETIME(6);default:CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6);comment:更新時間"`
}

func (p *Product) TableName() string {
	return "product"
}
