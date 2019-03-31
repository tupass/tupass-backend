// +build unit

package testing

import (
	"fmt"
	"testing"

	"github.com/tupass/tupass-backend/metric"
)

// TestCalculateLength tests function metric.CalculateLength().
func TestCalculateLength(t *testing.T) {
	testValues := []string{"hello", " a b c d e f ", "~#/", "aT1{", "", "kLKJiKJ8G6==JJ"}

	expectedOutput := []int{5, 13, 3, 4, 0, 14}
	t.Log("Testing CalculateLength() in length.go")
	for i := 0; i < len(expectedOutput); i++ {
		t.Logf("Testing: c: '%+q'", testValues[i])
		if test := metric.CalculateLength(testValues[i]); test != expectedOutput[i] {
			errMessage := fmt.Sprintf("output of CalculateLength('%+q') is not as expected. \n Result: %d \n Expected: %d", testValues[i], test, expectedOutput[i])
			t.Error(errMessage)
		}
	}
}
