package controller

import (
	"fmt"

	"shopping-cart/service/model/bo"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"
)

type ProductControllerInterface interface {
	List(ctx *gin.Context)
	Import(ctx *gin.Context)
}

type productController struct {
	in ctrlIn
}

func newProductController(in ctrlIn) ProductControllerInterface {
	return &productController{
		in: in,
	}
}

func (ctrl *productController) List(ctx *gin.Context) {

	product, err := ctrl.in.ProductCore.List(ctx)
	if err != nil {
		ctx.JSON(500, err.Error())
	}

	ctx.JSON(200, product)
}

func (ctrl *productController) Import(ctx *gin.Context) {

	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	xls, err := excelize.OpenReader(file)
	if err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	rows, err := xls.GetRows(xls.GetSheetName(xls.GetActiveSheetIndex()))
	if err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	products, err := ctrl.transfer(rows)
	if err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	if err := ctrl.in.ProductCore.Insert(ctx, products); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	ctx.JSON(200, "OK")
}

func (ctrl *productController) transfer(rows [][]string) ([]*bo.Product, error) {
	products := []*bo.Product{}
	for rIndex, row := range rows {
		if rIndex == 0 {
			continue
		}

		product := &bo.Product{}
		for cIndex, data := range row {
			// 產品名稱
			if cIndex == 0 {
				product.Name = data
			}
			// 產品成本
			if cIndex == 1 {
				cost, err := decimal.NewFromString(data)
				if err != nil {
					return nil, fmt.Errorf("set cost(%d_%d) at err:%v", rIndex, cIndex, err)
				}
				product.Cost = cost
			}

		}

		products = append(products, product)
	}

	return products, nil
}
