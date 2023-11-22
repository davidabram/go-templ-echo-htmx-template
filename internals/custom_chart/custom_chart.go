package custom_chart

import (
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/golang/freetype/truetype"
	"github.com/wcharczuk/go-chart/v2"
)

// Chart is what we're drawing.
type CustomChart struct {
	Title      string
	TitleStyle chart.Style

	ColorPalette chart.ColorPalette

	Width  int
	Height int
	DPI    float64

	Background chart.Style
	Canvas     chart.Style

	XAxis          chart.XAxis
	YAxis          chart.YAxis
	YAxisSecondary chart.YAxis

	Font        *truetype.Font
	defaultFont *truetype.Font

	Series   []CustomSeries
	Elements []CustomRenderable

	Log chart.Logger
}

// GetDPI returns the dpi for the chart.
func (c CustomChart) GetDPI(defaults ...float64) float64 {
	if c.DPI == 0 {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return chart.DefaultDPI
	}
	return c.DPI
}

// GetFont returns the text font.
func (c CustomChart) GetFont() *truetype.Font {
	if c.Font == nil {
		return c.defaultFont
	}
	return c.Font
}

// GetWidth returns the chart width or the default value.
func (c CustomChart) GetWidth() int {
	if c.Width == 0 {
		return chart.DefaultChartWidth
	}
	return c.Width
}

// GetHeight returns the chart height or the default value.
func (c CustomChart) GetHeight() int {
	if c.Height == 0 {
		return chart.DefaultChartHeight
	}
	return c.Height
}

// Render renders the chart with the given renderer to the given io.Writer.
func (c CustomChart) Render(rp CustomRendererProvider, w io.Writer) error {
	if len(c.Series) == 0 {
		return errors.New("please provide at least one series")
	}
	if err := c.checkHasVisibleSeries(); err != nil {
		return err
	}

	c.YAxisSecondary.AxisType = chart.YAxisSecondary

	r, err := rp(c.GetWidth(), c.GetHeight())
	if err != nil {
		return err
	}

	if c.Font == nil {
		defaultFont, err := chart.GetDefaultFont()
		if err != nil {
			return err
		}
		c.defaultFont = defaultFont
	}
	r.SetDPI(c.GetDPI(chart.DefaultDPI))

	c.drawBackground(r)

	var xt, yt, yta []chart.Tick
	xr, yr, yra := c.getRanges()
	canvasBox := c.getDefaultCanvasBox()
	xf, yf, yfa := c.getValueFormatters()

	chart.Debugf(c.Log, "chart; canvas box: %v", canvasBox)

	xr, yr, yra = c.setRangeDomains(canvasBox, xr, yr, yra)

	err = c.checkRanges(xr, yr, yra)
	if err != nil {
		r.Save(w)
		return err
	}

	if c.hasAxes() {
		xt, yt, yta = c.getAxesTicks(r, xr, yr, yra, xf, yf, yfa)
		canvasBox = c.getAxesAdjustedCanvasBox(r, canvasBox, xr, yr, yra, xt, yt, yta)
		xr, yr, yra = c.setRangeDomains(canvasBox, xr, yr, yra)

		chart.Debugf(c.Log, "chart; axes adjusted canvas box: %v", canvasBox)

		// do a second pass in case things haven't settled yet.
		xt, yt, yta = c.getAxesTicks(r, xr, yr, yra, xf, yf, yfa)
		canvasBox = c.getAxesAdjustedCanvasBox(r, canvasBox, xr, yr, yra, xt, yt, yta)
		xr, yr, yra = c.setRangeDomains(canvasBox, xr, yr, yra)
	}

	c.drawCanvas(r, canvasBox)
	c.drawAxes(r, canvasBox, xr, yr, yra, xt, yt, yta)
	for index, series := range c.Series {
		c.drawSeries(r, canvasBox, xr, yr, yra, series, index)
	}

	c.drawTitle(r)

	for _, a := range c.Elements {
		a(r, canvasBox, c.styleDefaultsElements())
	}

	return r.Save(w)
}

func (c CustomChart) checkHasVisibleSeries() error {
	var style chart.Style
	for _, s := range c.Series {
		style = s.GetStyle()
		if !style.Hidden {
			return nil
		}
	}
	return fmt.Errorf("chart render; must have (1) visible series")
}

func (c CustomChart) validateSeries() error {
	var err error
	for _, s := range c.Series {
		err = s.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c CustomChart) getRanges() (xrange, yrange, yrangeAlt chart.Range) {
	var minx, maxx float64 = math.MaxFloat64, -math.MaxFloat64
	var miny, maxy float64 = math.MaxFloat64, -math.MaxFloat64
	var minya, maxya float64 = math.MaxFloat64, -math.MaxFloat64

	seriesMappedToSecondaryAxis := false

	// note: a possible future optimization is to not scan the series values if
	// all axis are represented by either custom ticks or custom ranges.
	for _, s := range c.Series {
		if !s.GetStyle().Hidden {
			seriesAxis := s.GetYAxis()
			if bvp, isBoundedValuesProvider := s.(chart.BoundedValuesProvider); isBoundedValuesProvider {
				seriesLength := bvp.Len()
				for index := 0; index < seriesLength; index++ {
					vx, vy1, vy2 := bvp.GetBoundedValues(index)

					minx = math.Min(minx, vx)
					maxx = math.Max(maxx, vx)

					if seriesAxis == chart.YAxisPrimary {
						miny = math.Min(miny, vy1)
						miny = math.Min(miny, vy2)
						maxy = math.Max(maxy, vy1)
						maxy = math.Max(maxy, vy2)
					} else if seriesAxis == chart.YAxisSecondary {
						minya = math.Min(minya, vy1)
						minya = math.Min(minya, vy2)
						maxya = math.Max(maxya, vy1)
						maxya = math.Max(maxya, vy2)
						seriesMappedToSecondaryAxis = true
					}
				}
			} else if vp, isValuesProvider := s.(chart.ValuesProvider); isValuesProvider {
				seriesLength := vp.Len()
				for index := 0; index < seriesLength; index++ {
					vx, vy := vp.GetValues(index)

					minx = math.Min(minx, vx)
					maxx = math.Max(maxx, vx)

					if seriesAxis == chart.YAxisPrimary {
						miny = math.Min(miny, vy)
						maxy = math.Max(maxy, vy)
					} else if seriesAxis == chart.YAxisSecondary {
						minya = math.Min(minya, vy)
						maxya = math.Max(maxya, vy)
						seriesMappedToSecondaryAxis = true
					}
				}
			}
		}
	}

	if c.XAxis.Range == nil {
		xrange = &chart.ContinuousRange{}
	} else {
		xrange = c.XAxis.Range
	}

	if c.YAxis.Range == nil {
		yrange = &chart.ContinuousRange{}
	} else {
		yrange = c.YAxis.Range
	}

	if c.YAxisSecondary.Range == nil {
		yrangeAlt = &chart.ContinuousRange{}
	} else {
		yrangeAlt = c.YAxisSecondary.Range
	}

	if len(c.XAxis.Ticks) > 0 {
		tickMin, tickMax := math.MaxFloat64, -math.MaxFloat64
		for _, t := range c.XAxis.Ticks {
			tickMin = math.Min(tickMin, t.Value)
			tickMax = math.Max(tickMax, t.Value)
		}
		xrange.SetMin(tickMin)
		xrange.SetMax(tickMax)
	} else if xrange.IsZero() {
		xrange.SetMin(minx)
		xrange.SetMax(maxx)
	}

	if len(c.YAxis.Ticks) > 0 {
		tickMin, tickMax := math.MaxFloat64, -math.MaxFloat64
		for _, t := range c.YAxis.Ticks {
			tickMin = math.Min(tickMin, t.Value)
			tickMax = math.Max(tickMax, t.Value)
		}
		yrange.SetMin(tickMin)
		yrange.SetMax(tickMax)
	} else if yrange.IsZero() {
		yrange.SetMin(miny)
		yrange.SetMax(maxy)

		if !c.YAxis.Style.Hidden {
			delta := yrange.GetDelta()
			roundTo := chart.GetRoundToForDelta(delta)
			rmin, rmax := chart.RoundDown(yrange.GetMin(), roundTo), chart.RoundUp(yrange.GetMax(), roundTo)

			yrange.SetMin(rmin)
			yrange.SetMax(rmax)
		}
	}

	if len(c.YAxisSecondary.Ticks) > 0 {
		tickMin, tickMax := math.MaxFloat64, -math.MaxFloat64
		for _, t := range c.YAxis.Ticks {
			tickMin = math.Min(tickMin, t.Value)
			tickMax = math.Max(tickMax, t.Value)
		}
		yrangeAlt.SetMin(tickMin)
		yrangeAlt.SetMax(tickMax)
	} else if seriesMappedToSecondaryAxis && yrangeAlt.IsZero() {
		yrangeAlt.SetMin(minya)
		yrangeAlt.SetMax(maxya)

		if !c.YAxisSecondary.Style.Hidden {
			delta := yrangeAlt.GetDelta()
			roundTo := chart.GetRoundToForDelta(delta)
			rmin, rmax := chart.RoundDown(yrangeAlt.GetMin(), roundTo), chart.RoundUp(yrangeAlt.GetMax(), roundTo)
			yrangeAlt.SetMin(rmin)
			yrangeAlt.SetMax(rmax)
		}
	}

	return
}

func (c CustomChart) checkRanges(xr, yr, yra chart.Range) error {
	chart.Debugf(c.Log, "checking xrange: %v", xr)
	xDelta := xr.GetDelta()
	if math.IsInf(xDelta, 0) {
		return errors.New("infinite x-range delta")
	}
	if math.IsNaN(xDelta) {
		return errors.New("nan x-range delta")
	}
	if xDelta == 0 {
		return errors.New("zero x-range delta; there needs to be at least (2) values")
	}

	chart.Debugf(c.Log, "checking yrange: %v", yr)
	yDelta := yr.GetDelta()
	if math.IsInf(yDelta, 0) {
		return errors.New("infinite y-range delta")
	}
	if math.IsNaN(yDelta) {
		return errors.New("nan y-range delta")
	}

	if c.hasSecondarySeries() {
		chart.Debugf(c.Log, "checking secondary yrange: %v", yra)
		yraDelta := yra.GetDelta()
		if math.IsInf(yraDelta, 0) {
			return errors.New("infinite secondary y-range delta")
		}
		if math.IsNaN(yraDelta) {
			return errors.New("nan secondary y-range delta")
		}
	}

	return nil
}

func (c CustomChart) getDefaultCanvasBox() chart.Box {
	return c.Box()
}

func (c CustomChart) getValueFormatters() (x, y, ya chart.ValueFormatter) {
	for _, s := range c.Series {
		if vfp, isVfp := s.(chart.ValueFormatterProvider); isVfp {
			sx, sy := vfp.GetValueFormatters()
			if s.GetYAxis() == chart.YAxisPrimary {
				x = sx
				y = sy
			} else if s.GetYAxis() == chart.YAxisSecondary {
				x = sx
				ya = sy
			}
		}
	}
	if c.XAxis.ValueFormatter != nil {
		x = c.XAxis.GetValueFormatter()
	}
	if c.YAxis.ValueFormatter != nil {
		y = c.YAxis.GetValueFormatter()
	}
	if c.YAxisSecondary.ValueFormatter != nil {
		ya = c.YAxisSecondary.GetValueFormatter()
	}
	return
}

func (c CustomChart) hasAxes() bool {
	return !c.XAxis.Style.Hidden || !c.YAxis.Style.Hidden || !c.YAxisSecondary.Style.Hidden
}

func (c CustomChart) getAxesTicks(r chart.Renderer, xr, yr, yar chart.Range, xf, yf, yfa chart.ValueFormatter) (xticks, yticks, yticksAlt []chart.Tick) {
	if !c.XAxis.Style.Hidden {
		xticks = c.XAxis.GetTicks(r, xr, c.styleDefaultsAxes(), xf)
	}
	if !c.YAxis.Style.Hidden {
		yticks = c.YAxis.GetTicks(r, yr, c.styleDefaultsAxes(), yf)
	}
	if !c.YAxisSecondary.Style.Hidden {
		yticksAlt = c.YAxisSecondary.GetTicks(r, yar, c.styleDefaultsAxes(), yfa)
	}
	return
}

func (c CustomChart) getAxesAdjustedCanvasBox(r chart.Renderer, canvasBox chart.Box, xr, yr, yra chart.Range, xticks, yticks, yticksAlt []chart.Tick) chart.Box {
	axesOuterBox := canvasBox.Clone()
	if !c.XAxis.Style.Hidden {
		axesBounds := c.XAxis.Measure(r, canvasBox, xr, c.styleDefaultsAxes(), xticks)
		chart.Debugf(c.Log, "chart; x-axis measured %v", axesBounds)
		axesOuterBox = axesOuterBox.Grow(axesBounds)
	}
	if !c.YAxis.Style.Hidden {
		axesBounds := c.YAxis.Measure(r, canvasBox, yr, c.styleDefaultsAxes(), yticks)
		chart.Debugf(c.Log, "chart; y-axis measured %v", axesBounds)
		axesOuterBox = axesOuterBox.Grow(axesBounds)
	}
	if !c.YAxisSecondary.Style.Hidden && c.hasSecondarySeries() {
		axesBounds := c.YAxisSecondary.Measure(r, canvasBox, yra, c.styleDefaultsAxes(), yticksAlt)
		chart.Debugf(c.Log, "chart; y-axis secondary measured %v", axesBounds)
		axesOuterBox = axesOuterBox.Grow(axesBounds)
	}

	return canvasBox.OuterConstrain(c.Box(), axesOuterBox)
}

func (c CustomChart) setRangeDomains(canvasBox chart.Box, xr, yr, yra chart.Range) (chart.Range, chart.Range, chart.Range) {
	xr.SetDomain(canvasBox.Width())
	yr.SetDomain(canvasBox.Height())
	yra.SetDomain(canvasBox.Height())
	return xr, yr, yra
}

func (c CustomChart) hasSecondarySeries() bool {
	for _, s := range c.Series {
		if s.GetYAxis() == chart.YAxisSecondary {
			return true
		}
	}
	return false
}

func (c CustomChart) getBackgroundStyle() chart.Style {
	return c.Background.InheritFrom(c.styleDefaultsBackground())
}

func (c CustomChart) drawBackground(r chart.Renderer) {
	chart.Draw.Box(r, chart.Box{
		Right:  c.GetWidth(),
		Bottom: c.GetHeight(),
	}, c.getBackgroundStyle())
}

func (c CustomChart) getCanvasStyle() chart.Style {
	return c.Canvas.InheritFrom(c.styleDefaultsCanvas())
}

func (c CustomChart) drawCanvas(r chart.Renderer, canvasBox chart.Box) {
	chart.Draw.Box(r, canvasBox, c.getCanvasStyle())
}

func (c CustomChart) drawAxes(r chart.Renderer, canvasBox chart.Box, xrange, yrange, yrangeAlt chart.Range, xticks, yticks, yticksAlt []chart.Tick) {
	if !c.XAxis.Style.Hidden {
		c.XAxis.Render(r, canvasBox, xrange, c.styleDefaultsAxes(), xticks)
	}
	if !c.YAxis.Style.Hidden {
		c.YAxis.Render(r, canvasBox, yrange, c.styleDefaultsAxes(), yticks)
	}
	if !c.YAxisSecondary.Style.Hidden {
		c.YAxisSecondary.Render(r, canvasBox, yrangeAlt, c.styleDefaultsAxes(), yticksAlt)
	}
}

func (c CustomChart) drawSeries(r CustomRenderer, canvasBox chart.Box, xrange, yrange, yrangeAlt chart.Range, s CustomSeries, seriesIndex int) {
	if !s.GetStyle().Hidden {
		if s.GetYAxis() == chart.YAxisPrimary {
			s.Render(r, canvasBox, xrange, yrange, c.styleDefaultsSeries(seriesIndex))
		} else if s.GetYAxis() == chart.YAxisSecondary {
			s.Render(r, canvasBox, xrange, yrangeAlt, c.styleDefaultsSeries(seriesIndex))
		}
	}
}

func (c CustomChart) drawTitle(r chart.Renderer) {
	if len(c.Title) > 0 && !c.TitleStyle.Hidden {
		r.SetFont(c.TitleStyle.GetFont(c.GetFont()))
		r.SetFontColor(c.TitleStyle.GetFontColor(c.GetColorPalette().TextColor()))
		titleFontSize := c.TitleStyle.GetFontSize(chart.DefaultTitleFontSize)
		r.SetFontSize(titleFontSize)

		textBox := r.MeasureText(c.Title)

		textWidth := textBox.Width()
		textHeight := textBox.Height()

		titleX := (c.GetWidth() >> 1) - (textWidth >> 1)
		titleY := c.TitleStyle.Padding.GetTop(chart.DefaultTitleTop) + textHeight

		r.Text(c.Title, titleX, titleY)
	}
}

func (c CustomChart) styleDefaultsBackground() chart.Style {
	return chart.Style{
		FillColor:   c.GetColorPalette().BackgroundColor(),
		StrokeColor: c.GetColorPalette().BackgroundStrokeColor(),
		StrokeWidth: chart.DefaultBackgroundStrokeWidth,
	}
}

func (c CustomChart) styleDefaultsCanvas() chart.Style {
	return chart.Style{
		FillColor:   c.GetColorPalette().CanvasColor(),
		StrokeColor: c.GetColorPalette().CanvasStrokeColor(),
		StrokeWidth: chart.DefaultCanvasStrokeWidth,
	}
}

func (c CustomChart) styleDefaultsSeries(seriesIndex int) chart.Style {
	return chart.Style{
		DotColor:    c.GetColorPalette().GetSeriesColor(seriesIndex),
		StrokeColor: c.GetColorPalette().GetSeriesColor(seriesIndex),
		StrokeWidth: chart.DefaultSeriesLineWidth,
		Font:        c.GetFont(),
		FontSize:    chart.DefaultFontSize,
	}
}

func (c CustomChart) styleDefaultsAxes() chart.Style {
	return chart.Style{
		Font:        c.GetFont(),
		FontColor:   c.GetColorPalette().TextColor(),
		FontSize:    chart.DefaultAxisFontSize,
		StrokeColor: c.GetColorPalette().AxisStrokeColor(),
		StrokeWidth: chart.DefaultAxisLineWidth,
	}
}

func (c CustomChart) styleDefaultsElements() chart.Style {
	return chart.Style{
		Font: c.GetFont(),
	}
}

// GetColorPalette returns the color palette for the chart.
func (c CustomChart) GetColorPalette() chart.ColorPalette {
	if c.ColorPalette != nil {
		return c.ColorPalette
	}
	return chart.DefaultColorPalette
}

// Box returns the chart bounds as a box.
func (c CustomChart) Box() chart.Box {
	dpr := c.Background.Padding.GetRight(chart.DefaultBackgroundPadding.Right)
	dpb := c.Background.Padding.GetBottom(chart.DefaultBackgroundPadding.Bottom)

	return chart.Box{
		Top:    c.Background.Padding.GetTop(chart.DefaultBackgroundPadding.Top),
		Left:   c.Background.Padding.GetLeft(chart.DefaultBackgroundPadding.Left),
		Right:  c.GetWidth() - dpr,
		Bottom: c.GetHeight() - dpb,
	}
}
