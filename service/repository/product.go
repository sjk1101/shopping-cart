package repository

import (
	"context"

	"gorm.io/gorm"

	"shopping-cart/service/model/po"
)

type ProductRepositoryInterface interface {
	Get(ctx context.Context, db *gorm.DB, id int) (*po.Product, error)
}

type productRepository struct {
	in repositoryIn
}

func newProductRepository(in repositoryIn) ProductRepositoryInterface {
	return &productRepository{
		in: in,
	}
}

func (repo *productRepository) Get(ctx context.Context, db *gorm.DB, id int,
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
