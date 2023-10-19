package repository

import "gorm.io/gorm"

type CondFn func(tx *gorm.DB) *gorm.DB
