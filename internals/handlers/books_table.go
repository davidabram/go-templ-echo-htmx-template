
package handlers

import (
	"context"

	"davidabram/go-templ-echo-htmx-template/internals/templates"
	"davidabram/go-templ-echo-htmx-template/internals/db"

	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
)

func (a *App) BooksTable(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	page := &templates.Page{
		Title: "Books",
		Boosted: h.HxBoosted,
	}

	DB := &db.Database{}
	if err := DB.New(); err != nil {
		return err
	}

	books, err := DB.ListFullBookInfo()

	if err != nil {
		return err
	}

	components := templates.BooksTable(page, books)
	return components.Render(context.Background(), c.Response().Writer)
}
