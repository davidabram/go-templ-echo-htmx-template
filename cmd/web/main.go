package main

import (
	"github.com/labstack/echo/v4"
	"davidabram/go-templ-echo-htmx-template/internals/handlers"
)

func main() {
	e := echo.New()

	e.GET("/", handlers.Hello)

	e.Static("/", "dist")

	e.Logger.Fatal(e.Start(":1323"))
}
