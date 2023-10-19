package app

import (
	"sync"

	"shopping-cart/service/constant"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"

	"shopping-cart/service/binder"
	"shopping-cart/service/controller"
)

var (
	once            sync.Once
	shoppingCartApp *ShoppingCartApp
)

func InitShoppingCart(app ShoppingCartApp) {
	once.Do(func() {
		shoppingCartApp = &app
	},
	)
}

type ShoppingCartApp struct {
	dig.In

	BorCtrl     controller.BotControllerInterface
	ProductCtrl controller.ProductControllerInterface
	OrderCtrl   controller.OrderControllerInterface
}

func Run() {
	binder := binder.New()
	if err := binder.Invoke(InitShoppingCart); err != nil {
		panic(err)
	}

	engine := gin.New()
	setRoutes(engine)

	if err := engine.Run(constant.Address); err != nil {
		panic(err)
	}
}

func setRoutes(engine *gin.Engine) {
	setPrivateRoutes(engine) // ex: pprof
	setPublicRoutes(engine)
}

func setPublicRoutes(engine *gin.Engine) {
	engine.GET("products", shoppingCartApp.ProductCtrl.List)
	engine.POST("products/import", shoppingCartApp.ProductCtrl.Import)

	engine.POST("orders/import", shoppingCartApp.OrderCtrl.Import)
	engine.GET("orders/export", shoppingCartApp.OrderCtrl.Export)

	engine.POST("/callback", shoppingCartApp.BorCtrl.Repeat)
}

func setPrivateRoutes(engine *gin.Engine) {
	engine.GET("/ping", func(ctx *gin.Context) {
		ctx.String(200, "pong")
	})
}
