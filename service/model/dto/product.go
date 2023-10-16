package dto

import "github.com/shopspring/decimal"

type ProductResp struct {
	ProductID   string          `json:"productID"`
	ProductName string          `json:"productName"`
	Image       string          `json:"image"`
	Cost        decimal.Decimal `json:"price"`
	Quantity    int             `json:"quantity"`
}
