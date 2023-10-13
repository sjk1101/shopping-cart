package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"
	"shopping-cart/service/model/po"
	"time"
)

type OrderControllerInterface interface {
	Import(ctx *gin.Context)
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
		return
	}

	xls, err := excelize.OpenReader(file)
	if err != nil {
		return
	}
	rows, err := xls.GetRows(xls.GetSheetName(xls.GetActiveSheetIndex()))

	orders, err := ctrl.transfer(rows)
	if err != nil {
		panic(err)
	}
	if err := ctrl.in.OrderCore.Insert(ctx, orders); err != nil {
		panic(err)
	}

	ctx.JSON(200, "OK")
}

func (ctrl *orderController) transfer(rows [][]string) ([]*po.ShopeeCompletedOrder, error) {
	orderMap := make(map[string]*po.ShopeeCompletedOrder)
	for rIndex, row := range rows {
		if rIndex == 0 {
			continue
		}

		// 按Shopee給的excel格式（是訂單編號Ｘ品項）
		completeOrder := &po.ShopeeCompletedOrder{}
		for cIndex, data := range row {
			// 訂單編號
			if cIndex == 0 {
				completeOrder.OrderID = data
			}
			// 訂單狀態
			if cIndex == 1 {
				completeOrder.IsEstablished = data == "完成"
			}

			// 訂單成立日期
			if cIndex == 5 {
				// excel import time: 2023-09-11 10:23
				a := fmt.Sprintf("%s:00", data)
				date, err := time.Parse("2006-01-02 15:04:05", a)
				if err != nil {
					return nil, fmt.Errorf("set order created(%d_%d) at err:%v", rIndex, cIndex, err)
				}
				completeOrder.OrderCreatedAt = date
			}

			// 賣場優惠券
			if cIndex == 15 {
				couponDiscount, err := decimal.NewFromString(data)
				if err != nil {
					return nil, fmt.Errorf("set coupon discount(%d_%d) err:%v", rIndex, cIndex, err)
				}
				completeOrder.CouponDiscount = couponDiscount
			}

			// 成交手續費
			if cIndex == 18 {
				dealFee, err := decimal.NewFromString(data)
				if err != nil {
					return nil, fmt.Errorf("set deal fee(%d_%d)  err:%v", rIndex, cIndex, err)
				}
				completeOrder.DealFee = dealFee
			}

			// 活動服務費
			if cIndex == 19 {
				activityFee, err := decimal.NewFromString(data)
				if err != nil {
					return nil, fmt.Errorf("set activity fee(%d_%d) err:%v", rIndex, cIndex, err)
				}
				completeOrder.ActivityFee = activityFee
			}

			// 金流服務費
			if cIndex == 19 {
				cashFlowCost, err := decimal.NewFromString(data)
				if err != nil {
					return nil, fmt.Errorf("set cash flow cost(%d_%d) err:%v", rIndex, cIndex, err)
				}
				completeOrder.CashFlowCost = cashFlowCost
			}

			// 商品金額
			if cIndex == 26 {
				price, err := decimal.NewFromString(data)
				if err != nil {
					return nil, fmt.Errorf("set price(%d_%d) err:%v", rIndex, cIndex, err)
				}
				completeOrder.Price = completeOrder.Price.Add(price)
			}

			// 訂單完成才會有日期
			if completeOrder.IsEstablished {
				// 訂單完成日期
				if cIndex == 47 {
					a := fmt.Sprintf("%s:00", data)
					date, err := time.Parse("2006-01-02 15:04:05", a)
					if err != nil {
						return nil, fmt.Errorf("set order completed At(%d_%d) at err:%v", rIndex, cIndex, err)
					}
					completeOrder.OrderCompletedAt = &date

					// 撥款日
					allocate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
					completeOrder.AllocateAt = &allocate
				}
			}
		}

		// 是否已有此筆訂單編號紀錄
		obj, ok := orderMap[completeOrder.OrderID]
		if !ok {
			orderMap[completeOrder.OrderID] = completeOrder
		} else {
			obj.Price = obj.Price.Add(completeOrder.Price)
		}
	}

	completeOrders := []*po.ShopeeCompletedOrder{}
	for _, v := range orderMap {
		completeOrders = append(completeOrders, v)
	}

	return completeOrders, nil
}
