package custom_chart

import (
	"github.com/wcharczuk/go-chart/v2"
)

// Series is an alias to Renderable.
type CustomSeries interface {
	GetName() string
	GetYAxis() chart.YAxisType
	GetStyle() chart.Style
	Validate() error
	Render(r CustomRenderer, canvasBox chart.Box, xrange, yrange chart.Range, s chart.Style)
}
