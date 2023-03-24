package binder

import (
	"sync"

	"go.uber.org/dig"

	"shopping-cart/service/controller"
	"shopping-cart/service/repository"
)

var (
	binder *dig.Container
	once   sync.Once
)

func New() *dig.Container {
	once.Do(func() {
		binder = dig.New()

		// Controller
		if err := binder.Provide(controller.NewController); err != nil {
			panic(err)
		}

		// Repository
		if err := binder.Provide(repository.NewRepository); err != nil {
			panic(err)
		}
	})

	return binder
}
