package main

import (
	"GIN/routers"
	"GIN/settings"
)

func main() {
	router, setting := settings.Server()
	router.Static("static", "statics")
	router.LoadHTMLGlob("templates/*")

	routers.RouterUrl(router)

	err := setting.ListenAndServe()
	settings.ErrorFatal(err, "Ошибка сервера при включении")
}
