package custom_chart

import (
	"github.com/wcharczuk/go-chart/v2"
)

// Renderable is a function that can be called to render custom elements on the chart.
type CustomRenderable func(r CustomRenderer, canvasBox chart.Box, defaults chart.Style)
