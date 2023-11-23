package custom_chart

import (
	"fmt"
	"github.com/wcharczuk/go-chart/v2"
	"davidabram/go-templ-echo-htmx-template/internals/spline"
)

// Interface Assertions.
var (
	_ CustomSeries              = (*CustomContinuousSeries)(nil)
	_ chart.FirstValuesProvider = (*CustomContinuousSeries)(nil)
	_ chart.LastValuesProvider  = (*CustomContinuousSeries)(nil)
)

// CustomContinuousSeries represents a line on a chart.
type CustomContinuousSeries struct {
	Name  string
	Style chart.Style

	YAxis chart.YAxisType

	XValueFormatter chart.ValueFormatter
	YValueFormatter chart.ValueFormatter

	XValues []float64
	YValues []float64
}

// GetName returns the name of the time series.
func (css CustomContinuousSeries) GetName() string {
	return css.Name
}

// GetStyle returns the line style.
func (css CustomContinuousSeries) GetStyle() chart.Style {
	return css.Style
}

// Len returns the number of elements in the series.
func (css CustomContinuousSeries) Len() int {
	return len(css.XValues)
}

// GetValues gets the x,y values at a given index.
func (css CustomContinuousSeries) GetValues(index int) (float64, float64) {
	return css.XValues[index], css.YValues[index]
}

// GetFirstValues gets the first x,y values.
func (css CustomContinuousSeries) GetFirstValues() (float64, float64) {
	return css.XValues[0], css.YValues[0]
}

// GetLastValues gets the last x,y values.
func (css CustomContinuousSeries) GetLastValues() (float64, float64) {
	return css.XValues[len(css.XValues)-1], css.YValues[len(css.YValues)-1]
}

// GetValueFormatters returns value formatter defaults for the series.
func (css CustomContinuousSeries) GetValueFormatters() (x, y chart.ValueFormatter) {
	if css.XValueFormatter != nil {
		x = css.XValueFormatter
	} else {
		x = chart.FloatValueFormatter
	}
	if css.YValueFormatter != nil {
		y = css.YValueFormatter
	} else {
		y = chart.FloatValueFormatter
	}
	return
}

// GetYAxis returns which YAxis the series draws on.
func (css CustomContinuousSeries) GetYAxis() chart.YAxisType {
	return css.YAxis
}

// Render renders the series.
func (css CustomContinuousSeries) Render(r CustomRenderer, canvasBox chart.Box, xrange, yrange chart.Range, defaults chart.Style) {
	style := css.Style.InheritFrom(defaults)
	css.CustomLineSeries(r, canvasBox, xrange, yrange, style, css)
}

// Validate validates the series.
func (css CustomContinuousSeries) Validate() error {
	if len(css.XValues) == 0 {
		return fmt.Errorf("continuous series; must have xvalues set")
	}

	if len(css.YValues) == 0 {
		return fmt.Errorf("continuous series; must have yvalues set")
	}

	if len(css.XValues) != len(css.YValues) {
		return fmt.Errorf("continuous series; must have same length xvalues as yvalues")
	}
	return nil
}

func (css CustomContinuousSeries) CustomLineSeries(r CustomRenderer, canvasBox chart.Box, xrange, yrange chart.Range, style chart.Style, vs chart.ValuesProvider) {
	if vs.Len() == 0 {
		return
	}

	spline, err := spline.NewMonotoneSpline(css.XValues, css.YValues)

	if(err != nil) {
		panic(err)
	}

	cb := canvasBox.Bottom
	cl := canvasBox.Left

	v0x, v0y := vs.GetValues(0)
	x0 := cl + xrange.Translate(v0x)
	y0 := cb - yrange.Translate(v0y)

	yv0 := yrange.Translate(0)

	var vx, vy float64
	var x, y int

	if style.ShouldDrawStroke() && style.ShouldDrawFill() {
		style.GetFillOptions().WriteDrawingOptionsToRenderer(r)
		r.MoveTo(x0, y0)
		for i := 1; i < vs.Len(); i++ {
			vx, vy = vs.GetValues(i)
			x = cl + xrange.Translate(vx)
			y = cb - yrange.Translate(vy)
			r.LineTo(x, y)
		}
		r.LineTo(x, chart.MinInt(cb, cb-yv0))
		r.LineTo(x0, chart.MinInt(cb, cb-yv0))
		r.LineTo(x0, y0)
		r.Fill()
	}

	if style.ShouldDrawStroke() {
		style.GetStrokeOptions().WriteDrawingOptionsToRenderer(r)

		r.MoveTo(x0, y0)
		for i := 1; i < vs.Len(); i++ {
			vx, vy = vs.GetValues(i)
			x = cl + xrange.Translate(vx)
			y = cb - yrange.Translate(vy)

			vx1 := float64(vx)-0.8
			vx2 := float64(vx)-0.2

			vy1, err := spline.At(vx1)
			if(err != nil) {
				panic(err)
			}

			vy2, err := spline.At(vx2)
			if(err != nil) {
				panic(err)
			}

			cx1 := cl + xrange.Translate(vx1)
			cx2 := cl + xrange.Translate(vx2)
			cy1 := cb - yrange.Translate(vy1)
			cy2 := cb - yrange.Translate(vy2)

			fmt.Println(cx1, cy1, cx2, cy2)

			r.CubicCurveTo(cx1, cy1, cx2, cy2, x, y)
			x0 = x
			y0 = y
		}
		r.Stroke()
	}

	if style.ShouldDrawDot() {
		defaultDotWidth := style.GetDotWidth()

		style.GetDotOptions().WriteDrawingOptionsToRenderer(r)
		for i := 0; i < vs.Len(); i++ {
			vx, vy = vs.GetValues(i)
			x = cl + xrange.Translate(vx)
			y = cb - yrange.Translate(vy)

			dotWidth := defaultDotWidth
			if style.DotWidthProvider != nil {
				dotWidth = style.DotWidthProvider(xrange, yrange, i, vx, vy)
			}

			if style.DotColorProvider != nil {
				dotColor := style.DotColorProvider(xrange, yrange, i, vx, vy)

				r.SetFillColor(dotColor)
				r.SetStrokeColor(dotColor)
			}

			r.Circle(dotWidth, x, y)
			r.FillStroke()
		}
	}
}
