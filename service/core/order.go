package core

import (
	"context"
	"shopping-cart/service/model/po"
	"shopping-cart/service/thirdparty/database"
)

type OrderCoreInterface interface {
	Insert(ctx context.Context, orders []*po.ShopeeCompletedOrder) error
}

type orderCore struct {
	in coreIn
}

func newOrderCore(in coreIn) OrderCoreInterface {
	return &orderCore{
		in: in,
	}
}

func (c *orderCore) Insert(ctx context.Context, orders []*po.ShopeeCompletedOrder) error {

	db := database.Session()

	for index, v := range orders {
		v.ID = int64(index) + int64(1)
	}

	if err := c.in.OrderRepo.Create(ctx, db, orders); err != nil {
		panic(err)
	}
	return nil
}
