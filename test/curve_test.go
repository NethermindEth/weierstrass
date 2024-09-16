package weierstrass_test

import (
	"encoding/json"
	"math/big"
	"os"
	"testing"

	"github.com/mralj/weierstrass"
)

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

		if !curve.IsPointOnCurve(p) {
			t.Errorf("Point %d (%d, %d) is not on the curve", i, point.X, point.Y)
		}
	}

	for i, point := range curveData.PointsNotOnCurve {
		p := weierstrass.NewPoint(big.NewInt(point.X), big.NewInt(point.Y))

		if curve.IsPointOnCurve(p) {
			t.Errorf("Point %d (%d, %d) is on the curve", i, point.X, point.Y)
		}
	}
}

func TestPointAddition(t *testing.T) {
	jsonFile, err := os.Open("point_add_test.json")
	if err != nil {
		t.Fatalf("Failed to open JSON file: %v", err)
	}
	defer jsonFile.Close()

	var testData PointAddTestData
	decoder := json.NewDecoder(jsonFile)
	if err := decoder.Decode(&testData); err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}

	curve := weierstrass.NewCurve(
		big.NewInt(testData.A),
		big.NewInt(testData.B),
		big.NewInt(testData.P),
	)

	for i, test := range testData.Tests {
		P := weierstrass.NewPoint(big.NewInt(test.P.X), big.NewInt(test.P.Y))
		Q := weierstrass.NewPoint(big.NewInt(test.Q.X), big.NewInt(test.Q.Y))
		expectedR := weierstrass.NewPoint(big.NewInt(test.R.X), big.NewInt(test.R.Y))

		calculatedR := curve.AddPoints(P, Q)

		// Compare calculated R with expected R
		if !calculatedR.Eq(expectedR) {
			t.Errorf("Test case %d failed: P + Q != R", i)
			t.Errorf("P: %v", P)
			t.Errorf("Q: %v", Q)
			t.Errorf("Expected R: %v", expectedR)
			t.Errorf("Calculated R: %v", calculatedR)
		}
	}
}

func TestScalarMultiplication(t *testing.T) {
	jsonFile, err := os.Open("scalar_mul_test.json")
	if err != nil {
		t.Fatalf("Failed to open JSON file: %v", err)
	}
	defer jsonFile.Close()

	var testData ScalarMulTestData
	decoder := json.NewDecoder(jsonFile)
	if err := decoder.Decode(&testData); err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}

	// Create curve with parameters from JSON
	curve := weierstrass.NewCurve(
		big.NewInt(testData.A),
		big.NewInt(testData.B),
		big.NewInt(testData.P),
	)

	// Test each scalar multiplication case
	for i, test := range testData.Tests {
		P := weierstrass.NewPoint(big.NewInt(test.P.X), big.NewInt(test.P.Y))
		k := big.NewInt(test.K)
		expectedR := weierstrass.NewPoint(big.NewInt(test.R.X), big.NewInt(test.R.Y))

		// Perform scalar multiplication using your implementation
		calculatedR := curve.ScalarMulPoint(P, k)

		// Compare calculated R with expected R
		if !expectedR.Eq(calculatedR) {
			t.Errorf("Test case %d failed: k * P != R", i)
			t.Errorf("P: %v", P)
			t.Errorf("k: %v", k)
			t.Errorf("Expected R: %v", expectedR)
			t.Errorf("Calculated R: %v", calculatedR)
		}
	}
}

func TestScalarMulAddPoints(t *testing.T) {
	// Load test data from JSON file
	jsonFile, err := os.Open("scalar_mul_add_test.json")
	if err != nil {
		t.Fatalf("Failed to open JSON file: %v", err)
	}
	defer jsonFile.Close()

	var testData CombinedScalarTestData
	decoder := json.NewDecoder(jsonFile)
	if err := decoder.Decode(&testData); err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}

	// Create curve with parameters from JSON
	curve := weierstrass.NewCurve(
		big.NewInt(testData.A),
		big.NewInt(testData.B),
		big.NewInt(testData.P),
	)

	// Test each combined scalar multiplication and addition case
	for i, test := range testData.Tests {
		P := weierstrass.NewPoint(big.NewInt(test.P.X), big.NewInt(test.P.Y))
		Q := weierstrass.NewPoint(big.NewInt(test.Q.X), big.NewInt(test.Q.Y))
		expectedR := weierstrass.NewPoint(big.NewInt(test.R.X), big.NewInt(test.R.Y))
		k := big.NewInt(test.K)
		u := big.NewInt(test.U)

		// Perform combined scalar multiplication and addition using your implementation
		calculatedR := curve.ScalarMulAddPoints(P, Q, k, u)

		// Compare calculated R with expected R
		if !expectedR.Eq(calculatedR) {
			t.Errorf("Test case %d failed: kP + uQ != R", i)
			t.Errorf("P: %v", P)
			t.Errorf("k: %v", k)
			t.Errorf("Q: %v", Q)
			t.Errorf("u: %v", u)
			t.Errorf("Expected R: %v", expectedR)
			t.Errorf("Calculated R: %v", calculatedR)
		}
	}
}

type PointCoords struct {
	X int64 `json:"x"`
	Y int64 `json:"y"`
}

type CurveTestData struct {
	P                int64         `json:"p"`
	A                int64         `json:"a"`
	B                int64         `json:"b"`
	PointsOnCurve    []PointCoords `json:"points_on_curve"`
	PointsNotOnCurve []PointCoords `json:"points_not_on_curve"`
}

type PointAddTestData struct {
	P     int64     `json:"p"`
	A     int64     `json:"a"`
	B     int64     `json:"b"`
	Tests []AddTest `json:"tests"`
}

type AddTest struct {
	P PointCoords `json:"P"`
	Q PointCoords `json:"Q"`
	R PointCoords `json:"R"`
}

type ScalarMulTestData struct {
	P     int64           `json:"p"`
	A     int64           `json:"a"`
	B     int64           `json:"b"`
	Tests []ScalarMulTest `json:"tests"`
}

type ScalarMulTest struct {
	P PointCoords `json:"P"`
	K int64       `json:"k"`
	R PointCoords `json:"R"`
}

type CombinedScalarTestData struct {
	P     int64                `json:"p"`
	A     int64                `json:"a"`
	B     int64                `json:"b"`
	Tests []CombinedScalarTest `json:"tests"`
}

type CombinedScalarTest struct {
	P PointCoords `json:"P"`
	K int64       `json:"k"`
	Q PointCoords `json:"Q"`
	U int64       `json:"u"`
	R PointCoords `json:"R"`
}
