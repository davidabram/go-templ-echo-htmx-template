package handlers

import (
  "context"

  "github.com/labstack/echo/v4"
  "davidabram/go-templ-echo-htmx-template/internals/templates"
)

func Hello(c echo.Context) error {
  components := templates.Hello("David!")
  return components.Render(context.Background(), c.Response().Writer)
}
