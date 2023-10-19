package repository

import (
	"context"

	"gorm.io/gorm"

	"shopping-cart/service/model/po"
)

type ProductRepositoryInterface interface {
	Get(ctx context.Context, db *gorm.DB, id int) (*po.Product, error)
	Create(ctx context.Context, db *gorm.DB, products []*po.Product) error
	Find(ctx context.Context, db *gorm.DB, cond CondFn) ([]*po.Product, error)
}

type productRepository struct {
	in repositoryIn
}

func newProductRepository(in repositoryIn) ProductRepositoryInterface {
	return &productRepository{
		in: in,
	}
}

func (r *productRepository) Get(ctx context.Context, db *gorm.DB, id int,
) (*po.Product, error) {

	var product po.Product

	if err := db.
		Model(&po.Product{}).
		Where("id = ?", id).
		First(&product).
		Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) Create(ctx context.Context, db *gorm.DB, products []*po.Product) error {

	if err := db.Create(products).Error; err != nil {
		return err
	}

	return nil
}

func (r *productRepository) Find(ctx context.Context, db *gorm.DB, cond CondFn) ([]*po.Product, error) {

	res := []*po.Product{}
	if err := cond(db.Model(&po.Product{})).
		Find(&res).
		Error; err != nil {
		return nil, err
	}

	return res, nil
}
