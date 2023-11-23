package handlers

import (
	"context"

	"davidabram/go-templ-echo-htmx-template/internals/templates"


	"github.com/donseba/go-htmx"
	"github.com/labstack/echo/v4"
)

func (a *App) Charts(c echo.Context) error {
	r := c.Request()
	h := r.Context().Value(htmx.ContextRequestHeader).(htmx.HxRequestHeader)

	page := &templates.Page{
		Title:   "Charts",
		Boosted: h.HxBoosted,
	}

	var chartData []templates.TimeSeries;

	chartData = append(chartData, templates.TimeSeries{
		Title: "Series",
		XValues: []float64{1, 2, 3, 4, 5, 6, 7, 8},
		YValues: []float64{33, 24, 11, 23, 47, 122.4, 36, 10},
	})
	chartData = append(chartData, templates.TimeSeries{
		Title: "Series 2",
		XValues: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
		YValues: []float64{2, 0, 4, 3, 4, 0, 4, 3, 2},
	})

	chartData = append(chartData, templates.TimeSeries{
		Title: "Series 3",
		XValues: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
		YValues: []float64{10, 20, 20, 40, 50, 10, 5, 80, 80},
	})


	components := templates.Charts(page, chartData)
	return components.Render(context.Background(), c.Response().Writer)
}
