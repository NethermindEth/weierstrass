package weierstrass

import (
	"fmt"
	"math/big"
)

type Point struct {
	x, y *big.Int
}

func NewPoint(x, y *big.Int) Point {
	return Point{x, y}
}

func (p Point) String() string {
	return fmt.Sprintf("(%d, %d)", p.x, p.y)
}
