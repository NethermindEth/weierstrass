package weierstrass

import (
	"fmt"
	"math/big"
)

// Point represents affine point on Weierstrass curve
type Point struct {
	x, y *big.Int // coordinates
}

// Infinity is identity element on curve
var Infinity = Point{}

func (p Point) X() *big.Int {
	return p.x
}

func (p Point) Y() *big.Int {
	return p.y
}

func NewPoint(x, y *big.Int) Point {
	return Point{x, y}
}

func (p Point) Clone() Point {
	if p.Eq(Infinity) {
		return Point{}
	}

	return NewPoint(new(big.Int).Set(p.x), new(big.Int).Set(p.y))
}

func (p Point) String() string {
	return fmt.Sprintf("(%d, %d)", p.x, p.y)
}

// IsInf checks if point is Infinity
func (p Point) IsInf() bool {
	return p == Infinity
}

// Eq Checks if two points are Equal
func (p Point) Eq(q Point) bool {
	if p.IsInf() && q.IsInf() {
		return true
	}

	if p.IsInf() || q.IsInf() {
		return false
	}

	xIsEq := p.x.Cmp(q.x) == 0
	yIsEq := p.y.Cmp(q.y) == 0

	return xIsEq && yIsEq
}
