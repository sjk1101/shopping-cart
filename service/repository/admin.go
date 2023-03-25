package repository

import (
	"context"

	"gorm.io/gorm"

	"shopping-cart/service/model/po"
)

type AdminRepositoryInterface interface {
	Create(ctx context.Context, db *gorm.DB, user *po.AdminUser) error
}

type adminRepository struct {
	in repositoryIn
}

func newAdminRepository(in repositoryIn) AdminRepositoryInterface {
	return &adminRepository{
		in: in,
	}
}

func (repo *adminRepository) Create(
	ctx context.Context, db *gorm.DB, user *po.AdminUser) error {

	if err := db.Create(user).Error; err != nil {
		return err
	}

	return nil
}
