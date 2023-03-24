package dto

type ProductResp struct {
	ProductID   string `json:"productID"`
	ProductName string `json:"productName"`
	Image       string `json:"image"`
	Amount      int    `json:"amount"`
	Inventory   int    `json:"inventory"`
}
