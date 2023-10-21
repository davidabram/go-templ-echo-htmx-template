package handlers

import (
	"context"
	"encoding/json"

	"davidabram/go-templ-echo-htmx-template/internals/templates"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
)

func (a *App) Hello(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	b, _ := json.MarshalIndent(h, "", "\t")

	page := &templates.Page{
		Title: "OMG",
		Boosted: h.HxBoosted,
	}

	components := templates.Hello(page, "David!", string(b))
	return components.Render(context.Background(), c.Response().Writer)
}
