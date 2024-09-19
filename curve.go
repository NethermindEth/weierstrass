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

func (c Curve) ArePointCoordinatesWithinPrimeField(p *Point) bool {
	// invalid if x/y are bigger or eq to prime p
	xIsInvalid := p.x.Cmp(c.p) >= 0
	yIsInvalid := p.y.Cmp(c.p) >= 0

	if xIsInvalid || yIsInvalid {
		return false
	}

	xIsValid := p.x.Sign() >= 0
	yIsValid := p.y.Sign() >= 0

	return xIsValid && yIsValid
}

// IsOnCurve P(x,y) is on curve c if equation y^2 = x^3 + ax + b mod p holds
func (c Curve) IsOnCurve(p *Point) bool {
	left := new(big.Int).Mul(p.y, p.y) // y^2
	left.Mod(left, c.p)                // y^2 mod p

	right := new(big.Int).Exp(p.x, big.NewInt(3), nil) // x^3
	right.Add(right, new(big.Int).Mul(c.a, p.x))       // x^3 + ax
	right.Add(right, c.b)                              // x^3 + ax + b
	right.Mod(right, c.p)                              // (x^3 + ax + b) mod p

	return left.Cmp(right) == 0
}

// Double doubles point P on curve C
func (c Curve) Double(p *Point) *Point {
	numerator := new(big.Int).Mul(p.x, p.x) // x1^2
	numerator.Mul(numerator, big.NewInt(3)) // 3x1^2
	numerator.Add(numerator, c.a)           // 3x1^2 + a

	denominator := new(big.Int).Mul(p.y, big.NewInt(2)) // 2y1
	denominatorInv := new(big.Int).ModInverse(denominator, c.p)

	if denominatorInv == nil {
		// Division by zero, result is point at infinity.
		return Infinity
	}

	lambda := new(big.Int).Mul(numerator, denominatorInv) // (3x1^2 + a) / 2y1
	lambda.Mod(lambda, c.p)

	xr := new(big.Int).Mul(lambda, lambda)           // λ^2
	xr.Sub(xr, new(big.Int).Mul(p.x, big.NewInt(2))) // λ^2 - 2x1
	xr.Mod(xr, c.p)

	yr := new(big.Int).Sub(p.x, xr) // x1 - x3
	yr.Mul(lambda, yr)              // λ(x1 - x3)
	yr.Sub(yr, p.y)                 // λ(x1 - x3) - y1
	yr.Mod(yr, c.p)

	return NewPoint(xr, yr)
}

// Neg Per definition Negating point on Weierstrass curve is just setting its
// y coordinate to -y
func (c Curve) Neg(p *Point) *Point {
	negY := new(big.Int).Neg(p.y)

	// Modding since in finite Field Fp
	negY.Mod(negY, c.p)

	return NewPoint(p.x, negY)
}

// Add adds two points P and Q on curve C to get resulting point R
func (c Curve) Add(p, q *Point) *Point {
	// Infinity is Identity element, denoted by 0
	// Let P = 0, then P + Q = 0 + Q = 0
	if p.IsInf() {
		return q
	}

	if q.IsInf() {
		return p
	}

	// Let Q = -P, then P + Q = P + (-P) = P - P = 0
	if q.Eq(c.Neg(p)) {
		return Infinity
	}

	// Let Q = P, then P + Q = P + P = 2P
	if p.Eq(q) {
		return c.Double(p)
	}

	lambda := new(big.Int)
	numerator := new(big.Int).Sub(q.y, p.y)   // y2 - y1
	denominator := new(big.Int).Sub(q.x, p.x) // x2 - x1
	denominatorInv := new(big.Int).ModInverse(denominator, c.p)

	if denominatorInv == nil {
		// Division by zero, result is point at infinity.
		return Infinity
	}

	lambda.Mul(numerator, denominatorInv)
	lambda.Mod(lambda, c.p)

	xr := new(big.Int).Mul(lambda, lambda) // λ^2
	xr.Sub(xr, p.x)                        // λ^2 - x1
	xr.Sub(xr, q.x)                        // λ^2 - x1 - x2
	xr.Mod(xr, c.p)

	yr := new(big.Int).Sub(p.x, xr) // x1 - x3
	yr.Mul(lambda, yr)              // λ(x1 - x3)
	yr.Sub(yr, p.y)                 // λ(x1 - x3) - y1
	yr.Mod(yr, c.p)

	return NewPoint(xr, yr)
}

// ScalarMul multiplies point P on curve with some scalar k; kP
// Implemented using Montgomery ladder algorithm as described in
// https://www.matthieurivain.com/files/jcen11b.pdf (Algorithm 3, page 4)
//
// IMPORTANT!!!: works, but not prod ready, not only  not optimized
// But not side-channel attack resistant
func (c Curve) ScalarMul(p *Point, k *big.Int) *Point {
	r0 := Infinity.Clone()
	r1 := p.Clone()

	for i := k.BitLen() - 1; i >= 0; i-- {
		// R1−b ← R1−b + Rb
		// 2Rb
		if k.Bit(i) == 0 {
			r1 = c.Add(r1, r0)
			r0 = c.Double(r0)
		} else {
			r0 = c.Add(r0, r1)
			r1 = c.Double(r1)
		}
	}

	return r0
}

// ScalarMulAdd calculates R = kP + uQ
// Where P, Q and R are points on curve and k, u are scalars
//
// IMPORTANT!!!: works, but not prod ready, not only  not optimized
// But not side-channel attack resistant because ScalarMul isn't
func (c Curve) ScalarMulAdd(p, q *Point, k, u *big.Int) *Point {
	kp := c.ScalarMul(p, k)
	kq := c.ScalarMul(q, u)

	return c.Add(kp, kq)
}
