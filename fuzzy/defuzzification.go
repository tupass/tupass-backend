package fuzzy

import (
	"log"
	"math"
)

// Defuzzy returns the defuzzified value of given YArray (MF curve of Y (triangles))
// and YMFG (MF grade of Y determined using rules (the areas)) as float64.
// YArray and YMFG must be of same length, otherwise Defuzzy panics.
// If YMFG has zero degree (so the sum of all entries equal zero), Defuzzy panics.
// Defuzzy uses centroid function to calculate crisp output data.
func Defuzzy(YArray []float64, YMFG []float64) float64 {
	PanicIfUnequalLength(YArray, YMFG, "YArray", "YMFG")

	sum := 0.0
	for _, num := range YMFG {
		sum += num
	}

	zeroDegree := sum == 0
	if zeroDegree {
		log.Panicln("Error while trying to defuzzy: YMFG has zero degree")
	}

	return centroid(YArray, YMFG)
}

// centroid calculates the crisp output data and centroids for given YArray (MF curve of Y)
// and MF grades YMFG, thus returning the defuzzified result.
func centroid(YArray []float64, YMFG []float64) float64 {
	// small step-by-step square
	var StepSquare = 0.0
	var TotalSquare = 0.0
	epsilon := math.Nextafter(1, 2) - 1

	// for singleton case
	if len(YArray) == 1 {
		return YArray[0] * YMFG[0] / math.Max(YMFG[0], epsilon)
	}

	// other cases: sum of all delta of step-by-step area
	for i := 1; i < len(YArray); i++ {
		x1 := YArray[i-1]
		x2 := YArray[i]
		y1 := YMFG[i-1]
		y2 := YMFG[i]
		if !((y1 == 0.0 && y2 == 0.0) || x1 == x2) {
			var Delta, Square float64
			// rectangle
			if y1 == y2 {
				Delta = 0.5 * (x1 + x2)
				Square = (x2 - x1) * y1
				// triangle with height y2
			} else if y1 == 0.0 && y2 != 0.0 {
				Delta = 2.0/3.0*(x2-x1) + x1
				Square = 0.5 * (x2 - x1) * y2
				// triangle with height y1
			} else if y2 == 0.0 && y1 != 0.0 {
				Delta = 1.0/3.0*(x2-x1) + x1
				Square = 0.5 * (x2 - x1) * y1
			} else {
				Delta = (2.0/3.0*(x2-x1)*(y2+0.5*y1))/(y1+y2) + x1
				Square = 0.5 * (x2 - x1) * (y1 + y2)
			}
			StepSquare += Delta * Square
			TotalSquare += Square
		}
	}

	return StepSquare / math.Max(TotalSquare, epsilon)
}

// PanicIfUnequalLength panics if given arrays x and y are of unequal length.
// xName and yName should describe array x and y for debug purposes.
func PanicIfUnequalLength(x, y []float64, xName, yName string) {
	if len(x) != len(y) {
		log.Panicf("Expected %s and %s to have the same length, but found lengths of %d and %d instead.", xName, yName, len(x), len(y))
	}
}
