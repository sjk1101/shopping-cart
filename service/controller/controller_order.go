package controller

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"shopping-cart/service/model/bo"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"
)

type OrderControllerInterface interface {
	Import(ctx *gin.Context)
	Export(ctx *gin.Context)
}

type orderController struct {
	in ctrlIn
}

func newOrderController(in ctrlIn) OrderControllerInterface {
	return &orderController{
		in: in,
	}
}

func (ctrl *orderController) Import(ctx *gin.Context) {

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

	orders, err := ctrl.transfer(rows)
	if err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	if err := ctrl.in.OrderCore.Insert(ctx, orders); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	ctx.JSON(200, "OK")
}

func (ctrl *orderController) Export(ctx *gin.Context) {

	res, err := ctrl.in.OrderCore.Export(ctx)
	if err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	file, _ := res.WriteToBuffer()
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition",
		fmt.Sprintf(
			"attachment; filename=shopee_statistics_%s.xlsx",
			time.Now().UTC().Format("20060102150405")))

	ctx.Data(200, "application/octet-stream", file.Bytes())
}

func (ctrl *orderController) transfer(rows [][]string) ([]*bo.ShopeeOrderDetail, error) {
	details := []*bo.ShopeeOrderDetail{}
	for rIndex, row := range rows {
		if rIndex == 0 {
			continue
		}

		// 按Shopee給的excel格式（是訂單編號Ｘ品項）
		detail := &bo.ShopeeOrderDetail{}
		for cIndex, data := range row {
			// 訂單編號
			if cIndex == 0 {
				detail.OrderID = data
			}
			// 訂單狀態
			if cIndex == 1 {
				detail.IsEstablished = data == "完成"
			}

			// 訂單成立日期
			if cIndex == 5 {
				// excel import time: 2023-09-11 10:23
				a := fmt.Sprintf("%s:00", data)
				date, err := time.Parse("2006-01-02 15:04:05", a)
				if err != nil {
					return nil, fmt.Errorf("set order created(%d_%d) at err:%v", rIndex, cIndex, err)
				}
				detail.OrderCreatedAt = date
			}

			// 賣場優惠券
			if cIndex == 15 {
				couponDiscount, err := decimal.NewFromString(data)
				if err != nil {
					return nil, fmt.Errorf("set coupon discount(%d_%d) err:%v", rIndex, cIndex, err)
				}
				detail.CouponDiscount = couponDiscount
			}

			// 成交手續費
			if cIndex == 18 {
				dealFee, err := decimal.NewFromString(data)
				if err != nil {
					return nil, fmt.Errorf("set deal fee(%d_%d)  err:%v", rIndex, cIndex, err)
				}
				detail.DealFee = dealFee
			}

			// 活動服務費
			if cIndex == 19 {
				activityFee, err := decimal.NewFromString(data)
				if err != nil {
					return nil, fmt.Errorf("set activity fee(%d_%d) err:%v", rIndex, cIndex, err)
				}
				detail.ActivityFee = activityFee
			}

			// 金流服務費
			if cIndex == 20 {
				cashFlowCost, err := decimal.NewFromString(data)
				if err != nil {
					return nil, fmt.Errorf("set cash flow cost(%d_%d) err:%v", rIndex, cIndex, err)
				}
				detail.CashFlowCost = cashFlowCost
			}

			// 商品名稱
			if cIndex == 23 {
				detail.Product = strings.ReplaceAll(data, "🔥", "")
			}

			// 商品選項名稱
			if cIndex == 24 {
				detail.Product = detail.Product + "," + data
			}

			// 商品金額
			if cIndex == 26 {
				price, err := decimal.NewFromString(data)
				if err != nil {
					return nil, fmt.Errorf("set price(%d_%d) err:%v", rIndex, cIndex, err)
				}
				detail.Price = price
			}

			// 商品數量
			if cIndex == 29 {
				q, err := strconv.Atoi(data)
				if err != nil {
					return nil, fmt.Errorf("set quantity(%d_%d) err:%v", rIndex, cIndex, err)
				}
				detail.Quantity = q
			}

			// 訂單完成才會有日期
			if detail.IsEstablished {
				// 訂單完成日期
				if cIndex == 47 {
					a := fmt.Sprintf("%s:00", data)
					date, err := time.Parse("2006-01-02 15:04:05", a)
					if err != nil {
						return nil, fmt.Errorf("set order completed At(%d_%d) at err:%v", rIndex, cIndex, err)
					}
					detail.OrderCompletedAt = &date
				}
			}
		}

		details = append(details, detail)
	}

	return details, nil
}
