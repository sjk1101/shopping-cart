package repository

import (
	"sync"

	"go.uber.org/dig"
)

var (
	once sync.Once
	self *repository
)

func NewRepository(in repositoryIn) repositoryOut {
	once.Do(func() {
		self = &repository{
			in: in,
			out: repositoryOut{
				ProductRepo: newProductRepository(in),
				OrderRepo:   newOrderRepository(in),
			},
		}
	})

	return self.out
}

type repository struct {
	in  repositoryIn
	out repositoryOut
}

type repositoryIn struct {
	dig.In
}

type repositoryOut struct {
	dig.Out

	ProductRepo ProductRepositoryInterface
	OrderRepo   OrderRepositoryInterface
}
