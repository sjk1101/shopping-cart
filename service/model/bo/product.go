package bo

import (
	"github.com/shopspring/decimal"
)

type Product struct {
	Name     string
	Number   string
	Cost     decimal.Decimal
	Quantity int
	Image    string
}
