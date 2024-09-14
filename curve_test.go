package weierstrass_test

import (
	"math/big"
	"testing"

	"github.com/mralj/weierstrass"
)

func TestIsOnCurve(t *testing.T) {
	curve := weierstrass.NewCurve(big.NewInt(63), big.NewInt(60), big.NewInt(97))

	// Test a point on the curve
	x := big.NewInt(25)
	y := big.NewInt(24)
	p := weierstrass.NewPoint(x, y)

	if !curve.IsOnCurve(p) {
		t.Errorf("Point (%v, %v) should be on the curve", x, y)
	}

	// Test a point not on the curve
	x = big.NewInt(10)
	y = big.NewInt(20)
	p = weierstrass.NewPoint(x, y)

	if curve.IsOnCurve(p) {
		t.Errorf("Point (%v, %v) should not be on the curve", x, y)
	}
}
