package repository

import (
	"context"
	"gorm.io/gorm"
	"shopping-cart/service/model/po"
)

type OrderRepositoryInterface interface {
	Create(ctx context.Context, db *gorm.DB, orders []*po.ShopeeCompletedOrder) error
}

type orderRepository struct {
	in repositoryIn
}

func newOrderRepository(in repositoryIn) OrderRepositoryInterface {
	return &orderRepository{
		in: in,
	}
}

func (r *orderRepository) Create(ctx context.Context, db *gorm.DB, orders []*po.ShopeeCompletedOrder) error {

	if err := db.Create(orders).Error; err != nil {
		return err
	}

	return nil
}
