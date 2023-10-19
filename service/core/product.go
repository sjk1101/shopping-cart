package core

import (
	"context"

	"shopping-cart/service/model/bo"
	"shopping-cart/service/model/dto"
	"shopping-cart/service/model/po"
	"shopping-cart/service/thirdparty/database"

	"gorm.io/gorm"
)

type ProductCoreInterface interface {
	Insert(ctx context.Context, products []*bo.Product) error
	List(ctx context.Context) ([]*dto.ProductResp, error)
}

type productCore struct {
	in coreIn
}

func newProductCore(in coreIn) ProductCoreInterface {
	return &productCore{
		in: in,
	}
}

func (c *productCore) Insert(ctx context.Context, products []*bo.Product) error {

	db := database.Session()

	poProducts := []*po.Product{}
	for _, v := range products {
		poProducts = append(poProducts, &po.Product{
			Name:     v.Name,
			Number:   "",
			Cost:     v.Cost,
			Quantity: 0,
			Image:    "",
		})
	}

	if err := c.in.ProductRepo.Create(ctx, db, poProducts); err != nil {
		return err
	}

	return nil
}

func (c *productCore) List(ctx context.Context) ([]*dto.ProductResp, error) {

	db := database.Session()
	poProducts, err := c.in.ProductRepo.Find(ctx, db,
		func(tx *gorm.DB) *gorm.DB {
			return tx
		})
	if err != nil {
		return nil, err
	}

	res := []*dto.ProductResp{}
	for _, v := range poProducts {
		res = append(res, &dto.ProductResp{
			ProductID:   v.Number,
			ProductName: v.Name,
			Image:       v.Image,
			Cost:        v.Cost,
			Quantity:    v.Quantity,
		})
	}
	return res, nil
}
