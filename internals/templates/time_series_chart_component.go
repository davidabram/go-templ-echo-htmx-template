package templates

import (
	"bytes"
	"context"
	"io"

	"github.com/a-h/templ"
	"davidabram/go-templ-echo-htmx-template/internals/custom_chart"
	"github.com/wcharczuk/go-chart/v2"
)


type TimeSeries struct {
	Title string
	XValues []float64
	YValues []float64
}

func TimeSeriesChart(series TimeSeries) templ.Component {
	mainSeries := custom_chart.CustomContinuousSeries{
		Name:    series.Title,
		XValues: series.XValues,
		YValues: series.YValues,
		Style: chart.Style{
			DotWidth:  5,
		},
	}

	graph := custom_chart.CustomChart{
		Width: 1400,
		Height: 400,
		Series: []custom_chart.CustomSeries{
			mainSeries,
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(custom_chart.CustomSVG, buffer)

	if err != nil {
		panic(err)
	}

	html := buffer.String()


	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, html)
		return err
	})
}
