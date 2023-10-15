package po

import (
	"time"
)

type AdminUser struct {
	ID        int64      `gorm:"column:id;type:bigint;primary_key;NOT NULL;comment:id"`
	Account   string     `gorm:"column:account;type:varchar(20);index:idx_account;NOT NULL;comment:管理者account"`
	Name      string     `gorm:"column:name;type:varchar(100);NOT NULL;comment:管理者名稱"`
	Password  string     `gorm:"column:password;type:varchar(128);NOT NULL;comment:管理者密碼"`
	Salt      string     `gorm:"column:salt;type:varchar(36);NOT NULL;comment:加密金鑰"`
	Status    int        `gorm:"column:status;type:int;comment:管理者狀態"`
	CreatedAt *time.Time `gorm:"column:created_at;type:DATETIME(6);default:CURRENT_TIMESTAMP(6);comment:建立時間"`
	UpdatedAt *time.Time `gorm:"column:updated_at;type:DATETIME(6);default:CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6);comment:更新時間"`
}

func (u *AdminUser) TableName() string {
	return "admin_user"
}
