package spline

import (
	"sort"
	"fmt"
)

type MathError struct {
	Message string
}

func (e *MathError) Error() string {
	return fmt.Sprintf("Math error: %s", e.Message)
}

func tridiagonalMatrixThomas(a, b, c, d []float64) ([]float64, error) {
	n := len(b)
	if len(a)+1 != n || n != len(c)+1 || n != len(d) {
		return nil, &MathError{"invalid matrix size"}
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
