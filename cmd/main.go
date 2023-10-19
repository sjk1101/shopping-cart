package main

import (
	"shopping-cart/service/app"
	"shopping-cart/service/daemon"
	"shopping-cart/service/thirdparty/bot"
	"shopping-cart/service/thirdparty/database"
)

func main() {
	if err := database.Init(); err != nil {
		panic(err)
	}

	if err := bot.InitLineBot(); err != nil {
		panic(err)
	}

	daemon.Run()
	app.Run()
}
