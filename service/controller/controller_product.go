package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"shopping-cart/service/model/dto"
	"shopping-cart/service/thirdparty/database"
)

type ProductControllerInterface interface {
	Get(ctx *gin.Context)
}

type productController struct {
	in ctrlIn
}

func newProductController(in ctrlIn) ProductControllerInterface {
	return &productController{
		in: in,
	}
}

func (ctrl *productController) Get(ctx *gin.Context) {

	db := database.Session()

	product, err := ctrl.in.ProductRepo.Get(ctx, db, 1)
	if err != nil {
		 ctx.JSON(500, err.Error())
	}

	//mock
	resp := []dto.ProductResp{dto.ProductResp{
		ProductID:  fmt.Sprintf("%d",product.ID)  ,
		ProductName: product.Name,
		Image:       product.Image,
		Amount:      product.Amount,
		Inventory:   product.Inventory,
	}}

	ctx.JSON(200, resp)
}
