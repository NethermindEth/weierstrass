package weierstrass

import (
	"fmt"
	"math/big"
)

// Curve represents a Weierstrass curve.
// y^2 = x^3 + ax + b mod p
// p is modulus of prime field
type Curve struct {
	a, b, p *big.Int
}

func NewCurve(a, b, p *big.Int) Curve {
	return Curve{a, b, p}
}

func (c Curve) String() string {
	return fmt.Sprintf("y^2 = x^3 + %d*x + %d mod %d", c.a, c.b, c.p)
}

// Point P(x,y) is on curve c if equation y^2 = x^3 + ax + b mod p holds
func (c Curve) IsOnCurve(p Point) bool {
	left := new(big.Int).Mul(p.y, p.y) // y^2
	left.Mod(left, c.p)                // y^2 mod p

	right := new(big.Int).Exp(p.x, big.NewInt(3), nil) // x^3
	right.Add(right, new(big.Int).Mul(c.a, p.x))       // x^3 + ax
	right.Add(right, c.b)                              // x^3 + ax + b
	right.Mod(right, c.p)                              // (x^3 + ax + b) mod p

	return left.Cmp(right) == 0
}
