package main

import (
	"shopping-cart/service/app"
	"shopping-cart/service/thirdparty/database"
)

func main() {
	if err := database.Init(); err != nil {
		panic(err)
	}
	app.Run()
}
