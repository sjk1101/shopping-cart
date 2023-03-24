package app

import (
	"sync"

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
	ProductCtrl controller.ProductControllerInterface
}

func Run() {
	binder := binder.New()
	if err := binder.Invoke(InitShoppingCart); err != nil {
		panic(err)
	}

	engine := gin.New()
	setRoutes(engine)

	if err := engine.Run(); err != nil {
		panic(err)
	}
}

func setRoutes(engine *gin.Engine) {
	setPrivateRoutes(engine) // ex: pprof
	setPublicRoutes(engine)
}

func setPublicRoutes(engine *gin.Engine) {
	engine.GET("products", shoppingCartApp.ProductCtrl.Get)
}

func setPrivateRoutes(engine *gin.Engine) {
	engine.GET("/ping", func(ctx *gin.Context) {
		ctx.String(200, "pong")
	})
}
