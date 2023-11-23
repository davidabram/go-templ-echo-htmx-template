package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"davidabram/go-templ-echo-htmx-template/internals/templates"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
)

func (a *App) SiginIn(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	fmt.Println("hello, sign in", c.Get("isSignedIn"))
	fmt.Println("user", c.Get("clerk_user"))

	b, _ := json.MarshalIndent(h, "", "\t")

	page := &templates.Page{
		Title:   "Sign in",
		Boosted: h.HxBoosted,
	}

	components := templates.SignIn(page, "David!", string(b))
	return components.Render(context.Background(), c.Response().Writer)
}
