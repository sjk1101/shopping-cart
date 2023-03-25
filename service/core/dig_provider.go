package core

import (
	"sync"

	"go.uber.org/dig"

	"shopping-cart/service/repository"
)

var (
	once sync.Once
	self *core
)

func NewCore(in coreIn) coreOut {
	once.Do(func() {
		self = &core{
			in: in,
			out: coreOut{
				AdminCore: newAdminCore(in),
			},
		}
	})

	return self.out
}

type core struct {
	in  coreIn
	out coreOut
}

type coreIn struct {
	dig.In

	AdminRepo repository.AdminRepositoryInterface
}

type coreOut struct {
	dig.Out

	AdminCore AdminCoreInterface
}
