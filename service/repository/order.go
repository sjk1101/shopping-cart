package repository

import (
	"context"
	"strings"

	"shopping-cart/service/model/po"

	"gorm.io/gorm"
)

type OrderRepositoryInterface interface {
	Create(ctx context.Context, db *gorm.DB, orders []*po.ShopeeOrder) error
	CreateDetails(ctx context.Context, db *gorm.DB, details []*po.ShopeeOrderDetail) error
	GetShopeeStatistics(ctx context.Context, db *gorm.DB) ([]*po.ShopeeOrderStatistics, error)
	FindDetails(ctx context.Context, db *gorm.DB) ([]*po.ShopeeOrderDetail, error)
}

type orderRepository struct {
	in repositoryIn
}

func newOrderRepository(in repositoryIn) OrderRepositoryInterface {
	return &orderRepository{
		in: in,
	}
}

func (r *orderRepository) Create(ctx context.Context, db *gorm.DB, orders []*po.ShopeeOrder) error {

	if err := db.Create(orders).Error; err != nil {
		return err
	}

	return nil
}

func (r *orderRepository) CreateDetails(ctx context.Context, db *gorm.DB, details []*po.ShopeeOrderDetail) error {

	if err := db.Create(details).Error; err != nil {
		return err
	}

	return nil
}

func (r *orderRepository) GetShopeeStatistics(ctx context.Context, db *gorm.DB) ([]*po.ShopeeOrderStatistics, error) {

	cols := []string{
		"allocate_at",
		"SUM(total_product_price) AS total_product_price",
		"SUM(coupon_discount) AS coupon_discount",
		"SUM(deal_fee) AS deal_fee",
		"SUM(activity_fee) AS activity_fee",
		"SUM(cash_flow_cost) AS cash_flow_cost",
		"SUM(total_product_cost) AS total_product_cost",
		"SUM(total_product_price)-SUM(coupon_discount)-SUM(deal_fee)-SUM(activity_fee)-SUM(cash_flow_cost)-SUM(total_product_cost) AS net_income",
	}

	res := []*po.ShopeeOrderStatistics{}
	if err := db.
		Model(&po.ShopeeOrder{}).
		Select(strings.Join(cols, ", ")).
		Where("is_established IS true").
		Group("allocate_at").
		Scan(&res).
		Error; err != nil {
		return nil, err
	}

	return res, nil
}

func (r *orderRepository) FindDetails(ctx context.Context, db *gorm.DB,
) ([]*po.ShopeeOrderDetail, error) {

	details := []*po.ShopeeOrderDetail{}
	if err := db.Model(details).
		Find(&details).
		Order("order_created_at").
		Error; err != nil {
		return nil, err
	}

	return details, nil
}
