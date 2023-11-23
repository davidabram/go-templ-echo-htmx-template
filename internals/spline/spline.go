package spline

type Spline interface {
	At(x float64) (float64, error)
	Range(start, end, step float64) ([]float64, error)
}
