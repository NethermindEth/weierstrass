package weierstrass_test

import (
	"encoding/json"
	"math/big"
	"os"
	"testing"

	"github.com/mralj/weierstrass"
)

type CurveTestData struct {
	P                int64         `json:"p"`
	A                int64         `json:"a"`
	B                int64         `json:"b"`
	PointsOnCurve    []PointCoords `json:"points_on_curve"`
	PointsNotOnCurve []PointCoords `json:"points_not_on_curve"`
}

type PointCoords struct {
	X int64 `json:"x"`
	Y int64 `json:"y"`
}

func TestIsPointOnCurve(t *testing.T) {
	jsonFile, err := os.Open("curve_test_gen.json")
	if err != nil {
		t.Fatalf("Failed to open JSON file: %v", err)
	}
	defer jsonFile.Close()

	var curveData CurveTestData
	decoder := json.NewDecoder(jsonFile)
	if err := decoder.Decode(&curveData); err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}

	curve := weierstrass.NewCurve(
		big.NewInt(curveData.A),
		big.NewInt(curveData.B),
		big.NewInt(curveData.P),
	)

	for i, point := range curveData.PointsOnCurve {
		p := weierstrass.NewPoint(big.NewInt(point.X), big.NewInt(point.Y))

		if !curve.IsOnCurve(p) {
			t.Errorf("Point %d (%d, %d) is not on the curve", i, point.X, point.Y)
		}
	}

	for i, point := range curveData.PointsNotOnCurve {
		p := weierstrass.NewPoint(big.NewInt(point.X), big.NewInt(point.Y))

		if curve.IsOnCurve(p) {
			t.Errorf("Point %d (%d, %d) is on the curve", i, point.X, point.Y)
		}
	}
}
