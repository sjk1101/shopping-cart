package controller

import (
	"sync"

	"go.uber.org/dig"

	"shopping-cart/service/repository"
)

var (
	once sync.Once
	self *ctrl
)

func NewController(in ctrlIn) ctrlOut {
	once.Do(func() {
		self = &ctrl{
			in: in,
			out: ctrlOut{
				ProductCtrl: newProductController(in),
			},
		}
	})

	return self.out
}

type ctrl struct {
	in  ctrlIn
	out ctrlOut
}

type ctrlIn struct {
	dig.In

	ProductRepo repository.ProductRepositoryInterface
}

type ctrlOut struct {
	dig.Out
	ProductCtrl ProductControllerInterface
}
