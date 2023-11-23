package spline

import (
	"sort"
	"fmt"
)

type boundary uint

const (
	CubicFirstDeriv boundary = iota
	CubicSecondDeriv
)

type cubic struct {
	x []float64
	y []float64
	n int

	boundary
	f0 float64
	fn float64

	m    []float64
	segs []*cubicSegment
}

// Cx = (xr - x)(ar * (xr - x)^2 + br) + (x - xl)(al * (x - xl)^2 + bl)
type cubicSegment struct {
	xl float64
	xr float64
	al float64
	bl float64
	ar float64
	br float64
}

type CubicInterpolationError struct {
	Message string
}

func (e *CubicInterpolationError) Error() string {
	return fmt.Sprintf("Cubic interpolation error: %s", e.Message)
}


func tridiagonalMatrixThomas(a, b, c, d []float64) ([]float64, error) {
	n := len(b)
	if len(a)+1 != n || n != len(c)+1 || n != len(d) {
		return nil, &CubicInterpolationError{"invalid matrix size"}
	}

	d[0] /= b[0]
	if n == 1 {
		return d, nil
	}

	c[0] /= b[0]
	for i := 1; i < n-1; i++ {
		div := b[i] - a[i-1]*c[i-1]
		c[i] /= div
		d[i] = (d[i] - a[i-1]*d[i-1]) / div
	}
	d[n-1] = (d[n-1] - a[n-2]*d[n-2]) / (b[n-1] - a[n-2]*c[n-2])
	for i := n - 2; i >= 0; i-- {
		d[i] -= c[i] * d[i+1]
	}
	return d, nil
}
func findSegment(xs []float64, x float64) int {
	if x <= xs[0] {
		return 0
	}
	if l := len(xs); x >= xs[l-1] {
		return l - 2
	}
	return sort.Search(len(xs), func(i int) bool { return xs[i] > x }) - 1
}

func NewCubicSpline(x, y []float64) (Spline, error) {
	return NewNaturalCubicSpline(x, y, 0, 0)
}

//	f0: f''(x[0])
//	fn: f''(x[len(x)-1])
func NewNaturalCubicSpline(x, y []float64, f0, fn float64) (Spline, error){
	return newSpline(x, y, CubicSecondDeriv, f0, fn)
}

//	f0: f'(x[0])
//	fn: f'(x[len(x)-1])
func NewClampedCubicSpline(x, y []float64, f0, fn float64) (Spline, error) {
	return newSpline(x, y, CubicFirstDeriv, f0, fn)
}

func (c *cubic) At(x float64) (float64, error) {
	nSegs := c.n - 1
	if c.segs == nil {
		c.segs = make([]*cubicSegment, nSegs)
	}
	seg := findSegment(c.x, x)
	s := c.segs[seg]
	// if not populated
	if s == nil {
		if c.m == nil {
			err := c.calculateM()
			if err != nil {
				return 0.0, err
			}
		}
		h := c.x[seg+1] - c.x[seg]
		s = &cubicSegment{
			xl: c.x[seg],
			xr: c.x[seg+1],
			ar: c.m[seg] / 6 / h,
			al: c.m[seg+1] / 6 / h,
			br: (c.y[seg] - c.m[seg]*h*h/6) / h,
			bl: (c.y[seg+1] - c.m[seg+1]*h*h/6) / h,
		}
		c.segs[seg] = s
	}
	dxr := s.xr - x
	dxl := x - s.xl

	return dxr*(s.ar*dxr*dxr+s.br) + dxl*(s.al*dxl*dxl+s.bl), nil
}

func (c *cubic) Range(start, end, step float64) ([]float64, error) {
	if start > end {
		return nil, &CubicInterpolationError{"start must be smaller than end"}
	}
	n := int((end-start)/step) + 1
	v := make([]float64, n)
	x := start
	for i := 0; i < n; i++ {
		vi, err := c.At(x)
		v[i] = vi
		if err != nil {
			return nil, err
		}
		x += step
	}
	return v, nil
}

func newSpline(x, y []float64, b boundary, f0, fn float64) (Spline, error) {
	if len(x) != len(y) {
		return nil, &CubicInterpolationError{"array length mismatch"}
	}
	n := len(x)
	if !sort.Float64sAreSorted(x) {
		return nil, &CubicInterpolationError{"x must be in ascending order"}
	}
	xx := make([]float64, n)
	copy(xx, x)
	yy := make([]float64, n)
	copy(yy, y)
	return &cubic{
		x:        xx,
		y:        yy,
		n:        n,
		boundary: b,
		f0:       f0,
		fn:       fn,
	}, nil
}

func (c *cubic) calculateM() error {
	h := make([]float64, c.n)
	for i := 1; i < c.n; i++ {
		h[i] = c.x[i] - c.x[i-1]
	}

	mu := make([]float64, c.n)
	lambda := make([]float64, c.n)
	diag := make([]float64, c.n)
	d := make([]float64, c.n)
	for i := 1; i < c.n-1; i++ {
		mu[i] = h[i] / (h[i] + h[i+1])
		lambda[i] = 1 - mu[i]
		diag[i] = 2
		d[i] = 6 * (c.y[i-1]/h[i]/(h[i]+h[i+1]) - c.y[i]/h[i]/h[i+1] + c.y[i+1]/(h[i]+h[i+1])/h[i+1])
	}
	diag[0] = 2
	diag[c.n-1] = 2

	switch c.boundary {
	case CubicFirstDeriv:
		mu[c.n-1] = 1
		lambda[0] = 1
		d[0] = 6 * ((c.y[1]-c.y[0])/h[1] - c.f0) / h[1]
		d[c.n-1] = 6 * (c.fn - (c.y[c.n-1]-c.y[c.n-2])/h[c.n-1]) / h[c.n-1]
	case CubicSecondDeriv:
		d[0] = 2 * c.f0
		d[c.n-1] = 2 * c.fn
	}

	m, err := tridiagonalMatrixThomas(mu[1:], diag, lambda[:c.n-1], d)

	if err != nil {
		return err
	}

	c.m = m
	return nil
}
