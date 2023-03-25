package controller

import "github.com/gin-gonic/gin"

type AdminControllerInterface interface {
	Create(ctx *gin.Context)
}

type adminController struct {
	in ctrlIn
}

func newAdminController(in ctrlIn) AdminControllerInterface {
	return &adminController{
		in: in,
	}
}

func(ctrl adminController)Create(ctx *gin.Context){




}