package fuzzy

import (
	"errors"
	"log"
	"math"
	"sort"
)

// DetTriangleMF calculates triangular membership function (MF) for given float64 array inV and indices array LMR,
// containing the left, middle and right values of the triangle.
// It returns a float64 array of the triangular membership function values.
// In case LMR is malformatted, an error is returned.
func DetTriangleMF(inV []float64, LMR []int) ([]float64, error) {
	L, M, R := float64(LMR[0]), float64(LMR[1]), float64(LMR[2])
	if len(LMR) != 3 || L > M || M > R {
		return nil, errors.New("DetTriangleMF received invalid parameters: Either LMR did not have 3 elements or the indices were invalid")
	}

	MFVal := make([]float64, len(inV))
	var Pos int

	// determine values of MFVal on the left side of triangle
	if L != M {
		for Pos = 0; Pos < len(inV); Pos++ {
			if L < inV[Pos] && inV[Pos] < M {
				MFVal[Pos] = (inV[Pos] - L) / float64(M-L)
			}
		}
	}

	// determine values of MFVal on the right side of triangle
	if M != R {
		for Pos = 0; Pos < len(inV); Pos++ {
			if M < inV[Pos] && inV[Pos] < R {
				MFVal[Pos] = (R - inV[Pos]) / float64(R-M)
			}
		}
	}

	// find the position of middle M to set MFVal to 1
	for Pos = 0; Pos < len(inV); Pos++ {
		if inV[Pos] == M {
			break
		}
	}
	MFVal[Pos] = 1

	// for the very left of the graph we want to have all 1 on the left of M
	if L == M {
		for Pos = 0; Pos < len(inV); Pos++ {
			if inV[Pos] <= M {
				MFVal[Pos] = 1
			}
		}
	}

	// for the very right of the graph we want to have all 1 on the right of M
	if M == R {
		for Pos = 0; Pos < len(inV); Pos++ {
			if inV[Pos] >= M {
				MFVal[Pos] = 1
			}
		}
	}

	return MFVal, nil
}

// DetMFGrad returns the (discrete) membership function grade at given x-coordinate
// for given x-coordinates of datapoints XDP and y-coordinates of datapoints (MF grade) YDP.
// It returns an interpolated value that corresponds to the MF grade.
func DetMFGrad(XDP []float64, YDP []float64, X float64) float64 {
	return interpLinear1D(XDP, YDP, X)
}

// interpLinear1D returns the linear interpolant (not zero-extrapolated) to a 1-dimensional
// function with given discrete data points (xData, yData), evaluated at x-coordinate x.
// The sequence of x-coordinates must be sorted and the sequence of y-coordinates
// must be of the same length as xData, otherwise interpLinear1D panics.
// It returns the interpolated value as float64.
// Based on https://github.com/gonum/gonum/issues/328#issuecomment-357473654
func interpLinear1D(xData, yData []float64, x float64) (y float64) {
	PanicIfUnequalLength(xData, yData, "xData", "yData")

	if !sort.Float64sAreSorted(xData) {
		log.Panicln("Expected xData to be sorted, but it is not.")
	}

	// If xData already contains the to be interpolated x-coordinate x
	// at the beginning or end, return corresponding x-coordinate in yData.
	// Otherwise interpolate between the closest two points in xData near x.
	index := sort.SearchFloat64s(xData, x)
	switch index {
	case 0:
		y = yData[0]
	case len(xData):
		y = yData[len(yData)-1]
	default:
		y = interp1DBetweenPoints(
			x, xData[index-1], xData[index],
			yData[index-1], yData[index])
	}
	return
}

// interp1DBetweenPoints interpolates between two points (x0, y1) and (x1, y1).
func interp1DBetweenPoints(x, x0, x1, y0, y1 float64) float64 {
	return y0 + (x-x0)*(y1-y0)/(x1-x0)
}

// Arrange returns an array of evenly spaced (in steps of given step) float64 values within given interval [start, stop).
func Arrange(start, stop, step float64) []float64 {
	N := int(math.Ceil((stop - start) / step))
	rnge := make([]float64, N)
	for x := range rnge {
		rnge[x] = start + step*float64(x)
	}
	return rnge
}

//CalculateMembershipGradesForPredictability returns an float64 array of membership grades of given predictability
func CalculateMembershipGradesForPredictability(predictability float64) []float64 {
	// password predictability (P)
	P := Arrange(0, 100, .1)

	// define triangle ranges (H=hard, M=medium, E=easy)
	PHVar, PMVar, PEVar := []int{30, 30, 50}, []int{30, 50, 70}, []int{50, 70, 70}

	// input membership functions
	PH, _ := DetTriangleMF(P, PHVar)
	PM, _ := DetTriangleMF(P, PMVar)
	PE, _ := DetTriangleMF(P, PEVar)

	// determine membership grades
	var PList []float64
	PList = append(PList, DetMFGrad(P, PH, predictability))
	PList = append(PList, DetMFGrad(P, PM, predictability))
	PList = append(PList, DetMFGrad(P, PE, predictability))
	return PList
}

//CalculateMembershipGradesForLength returns a float64 array of the membership grades of given length
func CalculateMembershipGradesForLength(length float64) []float64 {
	// password length
	L := Arrange(0., 27., .1)

	// define triangle ranges (Vs=very short, S=short, M=medium, L=long, Vl=very long)
	LVsVar, LSVar, LMVar, LLVar, LVlVar := []int{2, 2, 6}, []int{4, 8, 12}, []int{10, 14, 18}, []int{16, 20, 24}, []int{22, 26, 26}

	// input membership functions
	LVs, _ := DetTriangleMF(L, LVsVar)
	LS, _ := DetTriangleMF(L, LSVar)
	LM, _ := DetTriangleMF(L, LMVar)
	LL, _ := DetTriangleMF(L, LLVar)
	LVl, _ := DetTriangleMF(L, LVlVar)

	// determine membership grades
	var LList []float64
	LList = append(LList, DetMFGrad(L, LVs, float64(length)))
	LList = append(LList, DetMFGrad(L, LS, float64(length)))
	LList = append(LList, DetMFGrad(L, LM, float64(length)))
	LList = append(LList, DetMFGrad(L, LL, float64(length)))
	LList = append(LList, DetMFGrad(L, LVl, float64(length)))
	return LList
}

//CalculateMembershipGradesForComplexity returns a float64 array of the membership grades of given complexity
func CalculateMembershipGradesForComplexity(complexity float64) []float64 {
	// password complexity
	C := Arrange(0., 680., .1)

	// define triangle ranges (Vs=very simple, S=simple, M=medium, C=complex, Vc=very complex)
	CVsVar, CSVar, CMVar, CCVar, CVcVar := []int{5, 5, 173}, []int{5, 173, 341}, []int{173, 341, 509}, []int{341, 509, 677}, []int{509, 677, 677}

	// input Membership functions
	CVs, _ := DetTriangleMF(C, CVsVar)
	CS, _ := DetTriangleMF(C, CSVar)
	CM, _ := DetTriangleMF(C, CMVar)
	CC, _ := DetTriangleMF(C, CCVar)
	CVc, _ := DetTriangleMF(C, CVcVar)

	// determine membership grades (membershipList5)
	var CList []float64
	CList = append(CList, DetMFGrad(C, CVs, complexity))
	CList = append(CList, DetMFGrad(C, CS, complexity))
	CList = append(CList, DetMFGrad(C, CM, complexity))
	CList = append(CList, DetMFGrad(C, CC, complexity))
	CList = append(CList, DetMFGrad(C, CVc, complexity))
	return CList
}
