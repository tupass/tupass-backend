// +build unit

package testing

import (
	"fmt"
	"testing"

	"github.com/tupass/tupass-backend/metric"
)

// TestComplexityReward	tests the function metric.ComplexityReward().
func TestComplexityReward(t *testing.T) {
	testValuesC := []int{0, 50, 53, 45, 2, 6}
	testValuesN := []int{0, 0, 46, 27, 7, 43}
	testValuesM := []int{0, 27, 0, 0, 0, 11}
	testValuesP := []int{0, 52, 0, 0, 0, 42}
	testValuesQ := []int{0, 2, 72, 16, 0, 3}

	expectedOutput := []float64{0, 37.5, 26.5, 22.5, 0.5, 6}

	t.Log("Testing metric.ComplexityReward()")

	for i := 0; i < len(expectedOutput); i++ {
		t.Logf("Testing: c: '%d', n: '%d', m: '%d', p: '%d', q: '%d'", testValuesC[i], testValuesN[i], testValuesM[i], testValuesP[i], testValuesQ[i])

		if test := metric.ComplexityReward(testValuesC[i], testValuesN[i], testValuesM[i], testValuesP[i], testValuesQ[i]); test != expectedOutput[i] {
			t.Error(errorMessageComplexityReward(testValuesC[i], testValuesN[i], testValuesM[i], testValuesP[i], testValuesQ[i], test, expectedOutput[i]))
		}
	}
}

// errorMessageComplexityReward returns a error message for given parameters
// for metric.ComplexityReward, actual and expected values.
func errorMessageComplexityReward(c int, n int, m int, p int, q int, output float64, expected float64) string {
	return fmt.Sprintf("output of ComplexityReward('%d', '%d', '%d', '%d', '%d') is not as expected. \n Result: %f \n Expected: %f", c, n, m, p, q, output, expected)
}

// TestCalculateComplexity tests the function metric.CalculateComplexity().
func TestCalculateComplexity(t *testing.T) {
	testValues := []string{"", "a", "5", "#", "passwort", "P@$$w0rt", " test ", " TEST", "Äoderä", "P55hj#"}

	expectedOutput := []float64{0.0, 6.5, 2.5, 8.25, 52, 213, 85, 68.5, 78, 131}
	t.Log("Testing metric.CalculateComplexity()")
	for i := 0; i < len(expectedOutput); i++ {
		t.Logf("Testing: string: '%s'", testValues[i])

		if test := metric.CalculateComplexity(testValues[i]); test != expectedOutput[i] {
			t.Error(errorMessageCalculateComplexity(testValues[i], test, expectedOutput[i]))
		}
	}
}

// errorMessageCalculateComplexity returns a error message for given password, actual and expected values when testing CalculateComplexity().
func errorMessageCalculateComplexity(input string, output float64, expected float64) string {
	return fmt.Sprintf("output of metric.CalculateComplexity('%s') is not as expected. \n Result: %f \n Expected: %f", input, output, expected)
}
