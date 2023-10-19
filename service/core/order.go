package core

import (
	"context"
	"fmt"
	"time"

	"shopping-cart/service/constant"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"shopping-cart/service/model/bo"
	"shopping-cart/service/model/po"
	"shopping-cart/service/thirdparty/database"

	"github.com/xuri/excelize/v2"
)

type OrderCoreInterface interface {
	Insert(ctx context.Context, orders []*bo.ShopeeOrderDetail) error
	Export(ctx context.Context) (*excelize.File, error)
}

type orderCore struct {
	in coreIn
}

func newOrderCore(in coreIn) OrderCoreInterface {
	return &orderCore{
		in: in,
	}
}

func (c *orderCore) Insert(ctx context.Context, details []*bo.ShopeeOrderDetail) error {

	db := database.Session()

	orderMap := make(map[string]*po.ShopeeOrder)
	orderDetails := []*po.ShopeeOrderDetail{}

	products, err := c.in.ProductRepo.Find(ctx, db,
		func(tx *gorm.DB) *gorm.DB {
			return tx
		})
	if err != nil {
		return err
	}

	productMap := make(map[string]*po.Product)
	for _, v := range products {
		productMap[v.Name] = v
	}

	for _, v := range details {
		product, ok := productMap[v.Product]
		if !ok {
			return fmt.Errorf("this product doesn't exists")
		}

		orderDetails = append(orderDetails, &po.ShopeeOrderDetail{
			OrderID:          v.OrderID,
			OrderCreatedAt:   v.OrderCreatedAt,
			IsEstablished:    v.IsEstablished,
			OrderCompletedAt: v.OrderCompletedAt,
			Product:          v.Product,
			Quantity:         v.Quantity,
			ProductPrice:     v.Price,
			ProductCost:      product.Cost,
			CouponDiscount:   v.CouponDiscount,
			DealFee:          v.DealFee,
			ActivityFee:      v.ActivityFee,
			CashFlowCost:     v.CashFlowCost,
		})

		// 計算總金額
		totalPrice := v.Price.Mul(decimal.NewFromInt(int64(v.Quantity)))
		// 計算總成本
		totalCost := product.Cost.Mul(decimal.NewFromInt(int64(v.Quantity)))

		// 撥款日
		var allocate *time.Time
		if v.IsEstablished {
			d := time.Date(v.OrderCompletedAt.Year(), v.OrderCompletedAt.Month(), v.OrderCompletedAt.Day(),
				0, 0, 0, 0, time.UTC)
			allocate = &d
		}

		if obj, ok := orderMap[v.OrderID]; !ok {
			orderMap[v.OrderID] = &po.ShopeeOrder{
				OrderID:           v.OrderID,
				OrderCreatedAt:    v.OrderCreatedAt,
				IsEstablished:     v.IsEstablished,
				OrderCompletedAt:  v.OrderCompletedAt,
				AllocateAt:        allocate,
				CouponDiscount:    v.CouponDiscount,
				DealFee:           v.DealFee,
				ActivityFee:       v.ActivityFee,
				CashFlowCost:      v.CashFlowCost,
				TotalProductPrice: totalPrice,
				TotalProductCost:  totalCost,
			}
		} else {
			obj.TotalProductPrice = obj.TotalProductPrice.Add(totalPrice)
			obj.TotalProductCost = obj.TotalProductCost.Add(totalCost)
		}
	}

	orders := []*po.ShopeeOrder{}
	for _, v := range orderMap {
		orders = append(orders, v)
	}

	f := func(tx *gorm.DB) error {
		if err := c.in.OrderRepo.CreateDetails(ctx, tx, orderDetails); err != nil {
			return err
		}

		if err := c.in.OrderRepo.Create(ctx, tx, orders); err != nil {
			return err
		}
		return nil
	}

	if err := db.Transaction(f); err != nil {
		return err
	}

	return nil
}

func (c *orderCore) Export(ctx context.Context) (*excelize.File, error) {

	file := excelize.NewFile()

	var err error
	file, err = c.setOrderSheet(ctx, file)
	if err != nil {
		return nil, err
	}

	file, err = c.setOrderDetailSheet(ctx, file)
	if err != nil {
		return nil, err
	}

	// default sheet1要放在最後delete才有效
	if err := file.DeleteSheet("Sheet1"); err != nil {
		return nil, err
	}

	return file, nil
}

func (c *orderCore) setOrderSheet(ctx context.Context, file *excelize.File) (*excelize.File, error) {
	sheetName := "order"
	index, err := file.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}

	if err := file.SetColWidth(sheetName, "A", "H", 20); err != nil {
		return nil, err
	}

	if err := file.SetCellValue(sheetName, "A1", "撥款日期"); err != nil {
		return nil, err
	}

	if err := file.SetCellValue(sheetName, "B1", "商品總金額"); err != nil {
		return nil, err
	}

	if err := file.SetCellValue(sheetName, "C1", "賣場優惠券折扣"); err != nil {
		return nil, err
	}

	if err := file.SetCellValue(sheetName, "D1", "成交手續費"); err != nil {
		return nil, err
	}

	if err := file.SetCellValue(sheetName, "E1", "活動服務費"); err != nil {
		return nil, err
	}

	if err := file.SetCellValue(sheetName, "F1", "金流服務費"); err != nil {
		return nil, err
	}

	if err := file.SetCellValue(sheetName, "G1", "商品總成本"); err != nil {
		return nil, err
	}

	if err := file.SetCellValue(sheetName, "H1", "淨利"); err != nil {
		return nil, err
	}

	// data
	db := database.Session()
	data, err := c.in.OrderRepo.GetShopeeStatistics(ctx, db)
	if err != nil {
		return nil, err
	}

	for i, v := range data {
		rowAt := i + 2
		if err := file.SetCellValue(sheetName, fmt.Sprintf("A%d", rowAt),
			fmt.Sprintf("%d-%d-%d", v.AllocateAt.Year(), v.AllocateAt.Month(), v.AllocateAt.Day())); err != nil {
			return nil, fmt.Errorf("set A row cell failed, err: %v", err)
		}

		if err := file.SetCellValue(sheetName, fmt.Sprintf("B%d", rowAt),
			v.TotalProductPrice); err != nil {
			return nil, fmt.Errorf("set B row cell failed, err: %v", err)
		}

		if err := file.SetCellValue(sheetName, fmt.Sprintf("C%d", rowAt),
			v.CouponDiscount); err != nil {
			return nil, fmt.Errorf("set C row cell failed, err: %v", err)
		}

		if err := file.SetCellValue(sheetName, fmt.Sprintf("D%d", rowAt),
			v.DealFee); err != nil {
			return nil, fmt.Errorf("set D row cell failed, err: %v", err)
		}

		if err := file.SetCellValue(sheetName, fmt.Sprintf("E%d", rowAt),
			v.ActivityFee); err != nil {
			return nil, fmt.Errorf("set E row cell failed, err: %v", err)
		}

		if err := file.SetCellValue(sheetName, fmt.Sprintf("F%d", rowAt),
			v.CashFlowCost); err != nil {
			return nil, fmt.Errorf("set F row cell failed, err: %v", err)
		}

		if err := file.SetCellValue(sheetName, fmt.Sprintf("G%d", rowAt),
			v.TotalProductCost); err != nil {
			return nil, fmt.Errorf("set G row cell failed, err: %v", err)
		}

		if err := file.SetCellValue(sheetName, fmt.Sprintf("H%d", rowAt),
			v.NetIncome); err != nil {
			return nil, fmt.Errorf("set H row cell failed, err: %v", err)
		}

		style1, _ := file.NewStyle(
			&excelize.Style{
				Alignment: &excelize.Alignment{
					Horizontal: "left",
					WrapText:   true,
				},
			},
		)

		if err := file.SetCellStyle(sheetName, fmt.Sprintf("A%d", rowAt),
			fmt.Sprintf("A%d", rowAt), style1); err != nil {
			return nil, fmt.Errorf("set row style failed, err: %v", err)
		}

		style2, _ := file.NewStyle(
			&excelize.Style{
				Alignment: &excelize.Alignment{
					Horizontal: "right",
					WrapText:   true,
				},
			},
		)

		if err := file.SetCellStyle(sheetName, fmt.Sprintf("B%d", rowAt),
			fmt.Sprintf("H%d", rowAt), style2); err != nil {
			return nil, fmt.Errorf("set row style failed, err: %v", err)
		}
	}

	file.SetActiveSheet(index)

	return file, nil
}

func (c *orderCore) setOrderDetailSheet(ctx context.Context, file *excelize.File) (*excelize.File, error) {
	sheetName := "order_detail"
	index, err := file.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}

	if err := file.SetColWidth(sheetName, "A", "L", 20); err != nil {
		return nil, err
	}

	if err := file.SetCellValue(sheetName, "A1", "訂單編號"); err != nil {
		return nil, err
	}

	if err := file.SetCellValue(sheetName, "B1", "訂單建立時間"); err != nil {
		return nil, err
	}

	if err := file.SetCellValue(sheetName, "C1", "訂單是否建立"); err != nil {
		return nil, err
	}

	if err := file.SetCellValue(sheetName, "D1", "訂單完成時間"); err != nil {
		return nil, err
	}

	if err := file.SetCellValue(sheetName, "E1", "商品名稱"); err != nil {
		return nil, err
	}

	if err := file.SetCellValue(sheetName, "F1", "商品價格"); err != nil {
		return nil, err
	}

	if err := file.SetCellValue(sheetName, "G1", "商品成本"); err != nil {
		return nil, err
	}

	if err := file.SetCellValue(sheetName, "H1", "商品數量"); err != nil {
		return nil, err
	}

	if err := file.SetCellValue(sheetName, "I1", "賣場優惠券折扣"); err != nil {
		return nil, err
	}

	if err := file.SetCellValue(sheetName, "J1", "成交手續費"); err != nil {
		return nil, err
	}

	if err := file.SetCellValue(sheetName, "K1", "活動服務費"); err != nil {
		return nil, err
	}

	if err := file.SetCellValue(sheetName, "L1", "金流服務費"); err != nil {
		return nil, err
	}

	// data
	db := database.Session()
	data, err := c.in.OrderRepo.FindDetails(ctx, db)
	if err != nil {
		return nil, err
	}

	currentOrderID := ""
	orderCount := 0
	for i, v := range data {
		var (
			rowAt          = i + 2
			couponDiscount = decimal.Zero
			dealFee        = decimal.Zero
			activityFee    = decimal.Zero
			cashFlowCost   = decimal.Zero
			fill           = excelize.Fill{}
		)

		if currentOrderID != v.OrderID {
			orderCount++

			couponDiscount = v.CouponDiscount
			dealFee = v.DealFee
			activityFee = v.ActivityFee
			cashFlowCost = v.CashFlowCost

			if orderCount%2 == 0 {
				fill = excelize.Fill{
					Type:    "pattern",
					Pattern: 1,
					Color:   []string{"#C7C6C1"},
					Shading: 0,
				}
			}
			currentOrderID = v.OrderID
		}

		if err := file.SetCellValue(sheetName, fmt.Sprintf("A%d", rowAt),
			v.OrderID); err != nil {
			return nil, fmt.Errorf("set A row cell failed, err: %v", err)
		}

		if err := file.SetCellValue(sheetName, fmt.Sprintf("B%d", rowAt),
			v.OrderCreatedAt.Format(constant.DateTimeFormat)); err != nil {
			return nil, fmt.Errorf("set B row cell failed, err: %v", err)
		}

		if err := file.SetCellValue(sheetName, fmt.Sprintf("C%d", rowAt),
			v.IsEstablished); err != nil {
			return nil, fmt.Errorf("set C row cell failed, err: %v", err)
		}

		completedAt := ""
		if v.IsEstablished {
			completedAt = v.OrderCompletedAt.Format(constant.DateTimeFormat)
		}

		if err := file.SetCellValue(sheetName, fmt.Sprintf("D%d", rowAt),
			completedAt); err != nil {
			return nil, fmt.Errorf("set D row cell failed, err: %v", err)
		}

		if err := file.SetCellValue(sheetName, fmt.Sprintf("E%d", rowAt),
			v.Product); err != nil {
			return nil, fmt.Errorf("set E row cell failed, err: %v", err)
		}

		if err := file.SetCellValue(sheetName, fmt.Sprintf("F%d", rowAt),
			v.ProductPrice); err != nil {
			return nil, fmt.Errorf("set F row cell failed, err: %v", err)
		}

		if err := file.SetCellValue(sheetName, fmt.Sprintf("G%d", rowAt),
			v.ProductCost); err != nil {
			return nil, fmt.Errorf("set G row cell failed, err: %v", err)
		}

		if err := file.SetCellValue(sheetName, fmt.Sprintf("H%d", rowAt),
			v.Quantity); err != nil {
			return nil, fmt.Errorf("set H row cell failed, err: %v", err)
		}

		if err := file.SetCellValue(sheetName, fmt.Sprintf("I%d", rowAt),
			couponDiscount); err != nil {
			return nil, fmt.Errorf("set I row cell failed, err: %v", err)
		}
		if err := file.SetCellValue(sheetName, fmt.Sprintf("J%d", rowAt),
			dealFee); err != nil {
			return nil, fmt.Errorf("set J row cell failed, err: %v", err)
		}

		if err := file.SetCellValue(sheetName, fmt.Sprintf("K%d", rowAt),
			activityFee); err != nil {
			return nil, fmt.Errorf("set K row cell failed, err: %v", err)
		}

		if err := file.SetCellValue(sheetName, fmt.Sprintf("L%d", rowAt),
			cashFlowCost); err != nil {
			return nil, fmt.Errorf("set L row cell failed, err: %v", err)
		}

		style1, _ := file.NewStyle(
			&excelize.Style{
				Alignment: &excelize.Alignment{
					Horizontal: "left",
					WrapText:   true,
				},
				Fill: fill,
			},
		)

		if err := file.SetCellStyle(sheetName, fmt.Sprintf("A%d", rowAt),
			fmt.Sprintf("E%d", rowAt), style1); err != nil {
			return nil, fmt.Errorf("set row style failed, err: %v", err)
		}

		style2, _ := file.NewStyle(
			&excelize.Style{
				Alignment: &excelize.Alignment{
					Horizontal: "right",
					WrapText:   true,
				},
				Fill: fill,
			},
		)

		if err := file.SetCellStyle(sheetName, fmt.Sprintf("F%d", rowAt),
			fmt.Sprintf("L%d", rowAt), style2); err != nil {
			return nil, fmt.Errorf("set row style failed, err: %v", err)
		}
	}

	file.SetActiveSheet(index)

	return file, nil
}
