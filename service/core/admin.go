package core

import "github.com/gin-gonic/gin"

type AdminCoreInterface interface {
	Create(ctx *gin.Context)
}

type adminCore struct {
	in coreIn
}

func newAdminCore(in coreIn) AdminCoreInterface {
	return &adminCore{
		in: in,
	}
}

func(core adminCore)Create(ctx *gin.Context){




}