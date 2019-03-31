package fes

import (
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/tupass/tupass-backend/fuzzy"
)

// password strength
var s = fuzzy.Arrange(0, 100, .1)

// define triangle ranges (Vw=very weak, W=weak, M=medium, S=strong, Vs=very strong)
var sVwVar, sWVar, sMVar, sSVar, sVsVar = []int{10, 10, 20}, []int{20, 30, 40}, []int{40, 50, 60}, []int{60, 70, 80}, []int{80, 90, 90}

// ruleDict is stores premises (l, c, p) and conclusion (s)
// l(enght) = (very short, short, medium, long, very long) = (0, 1, 2, 3, 4)
// c(omplexity) = (very simple, simple, medium, complex, very complex) = (0, 1, 2, 3, 4)
// p(redictability) = (hard, medium, easy) = (0, 1, 2)
// s(trength) = (very weak, weak, medium, strong, very strong) = (0, 1, 2, 3, 4)
var ruleDict = map[string]int{
	"4,4,0": 4, // 1
	"4,3,0": 4, // 2
	"4,2,0": 4, // 3
	"4,1,0": 3, // 4
	"4,0,0": 3, // 5
	"4,4,1": 3, // 6
	"4,3,1": 3, // 7
	"4,2,1": 3, // 8
	"4,1,1": 2, // 9
	"4,0,1": 2, // 10
	"4,4,2": 2, // 11
	"4,3,2": 2, // 12
	"4,2,2": 2, // 13
	"4,1,2": 1, // 14
	"4,0,2": 1, // 15
	"3,4,0": 4, // 16
	"3,3,0": 4, // 17
	"3,2,0": 4, // 18
	"3,1,0": 3, // 19
	"3,0,0": 2, // 20
	"3,4,1": 3, // 21
	"3,3,1": 3, // 22
	"3,2,1": 3, // 23
	"3,1,1": 2, // 24 //changed in V2
	"3,0,1": 1, // 25
	"3,*,2": 1, // 26
	"2,4,0": 3, // 27
	"2,3,0": 3, // 28
	"2,2,0": 3, // 29
	"2,1,0": 2, // 30
	"2,0,0": 1, // 31
	"2,4,1": 2, // 32
	"2,3,1": 2, // 33
	"2,2,1": 2, // 34 //changed in V2
	"2,1,1": 2, // 35 //changed in V2
	"2,0,1": 1, // 36
	"2,*,2": 1, // 37
	"1,4,0": 2, // 38
	"1,3,0": 2, // 39
	"1,2,0": 2, // 40
	"1,1,0": 2, // 41 //changed in V2
	"1,0,0": 0, // 42
	"1,4,1": 2, // 43 //changed in V2
	"1,3,1": 2, // 44 //changed in V2
	"1,2,1": 2, // 45 //changed in V2
	"1,1,1": 1, // 46 //changed in V2
	"1,0,1": 0, // 47
	"1,*,2": 0, // 48
	"0,*,*": 0} // 49

// updateOutputDict updates the output grade dict for given name of the strength set and
// Background: if there is only 1 rule can be applied, the membership grade of the output to the set will be the
// 	result of the rule, but if there are more than one rules, the "maximum of the membership grade" strategy will be
// 	used  => take the maximum of the results for this set
//     @param value: the result of the considering rule
func updateOutputDict(setName int, value float64, outputGradesDict *[5]float64) {

	//calculate and append output grades
	if outputGradesDict[setName] != 0.0 {
		//if there are more than one rules which can be applied, take the maximum result
		(*outputGradesDict)[setName] = math.Max(value, outputGradesDict[setName])
	} else {
		//if there is still no result for this set
		(*outputGradesDict)[setName] = value
	}
}

// Calculation for each rule in inference engine
// pos1: position for length index
// pos2: position for complexity index
// pos3: position for predictability index
// gradeL: membership grade of length input
// gradeC: membership grade of complexity input
// gradeP: membership grade of predictability input
func calculateForEachRule(pos1, pos2, pos3 string, gradeL, gradeC, gradeP float64, outputGradesDict *[5]float64) {
	//rule = '0,*,*'
	rule := strings.Join([]string{pos1, pos2, pos3}, ",")
	//get the name of the strength set
	strengthSet := ruleDict[rule]

	//calculate the grade result of output for the rule
	if gradeL == -1 {
		print("No rules with L = * in the rule set")
	} else {
		if gradeC == -1 { //Rule 26, 37, 48, 49
			if gradeP == -1 { //Rule 49
				ruleResult := gradeL
				updateOutputDict(strengthSet, ruleResult, outputGradesDict)
			} else { //Rule 26, 37, 48
				ruleResult := math.Min(gradeL, gradeP)
				updateOutputDict(strengthSet, ruleResult, outputGradesDict)
			}
		} else {
			if gradeP == -1 {
				log.Printf("No rules with L=%s and C=%s but P = *", pos1, pos2)
			} else { //Other normal rules
				ruleResult := math.Min(gradeL, math.Min(gradeC, gradeP))
				updateOutputDict(strengthSet, ruleResult, outputGradesDict)
			}
		}
	}
}

// GetStrengthByMembershipGrades returns the total strength based on given membership grades for length, complexity and predicatbility
func GetStrengthByMembershipGrades(LList, CList, PList []float64) float64 {
	var outputGradesDict = [5]float64{0.0, 0.0, 0.0, 0.0, 0.0}
	for indexL, itemL := range LList {
		if itemL == 0 {
			continue
		} else {
			for indexC, itemC := range CList {
				if itemC == 0 {
					continue
				} else {
					for indexP, itemP := range PList {
						if itemP == 0 {
							continue
						} else {
							if indexL == 0 { //Rule 49
								calculateForEachRule(strconv.Itoa(indexL), "*", "*", itemL, -1, -1, &outputGradesDict)
							} else if (indexL == 1) && (indexP == 2) { //Rule 48
								//rule = '1,*,2'
								calculateForEachRule(strconv.Itoa(indexL), "*", strconv.Itoa(indexP), itemL, -1, itemP, &outputGradesDict)
							} else if (indexL == 2) && (indexP == 2) { //Rule 37
								//rule = '2,*,2'
								calculateForEachRule(strconv.Itoa(indexL), "*", strconv.Itoa(indexP), itemL, -1, itemP, &outputGradesDict)
							} else if (indexL == 3) && (indexP == 2) { //Rule 26
								//rule = '3,*,2'
								calculateForEachRule(strconv.Itoa(indexL), "*", strconv.Itoa(indexP), itemL, -1, itemP, &outputGradesDict)
							} else { //other rules
								calculateForEachRule(strconv.Itoa(indexL), strconv.Itoa(indexC), strconv.Itoa(indexP), itemL, itemC, itemP, &outputGradesDict)
							}
						}
					}
				}
			}
		}
	}

	//output membership functions:
	sVw, _ := fuzzy.DetTriangleMF(s, sVwVar)
	sW, _ := fuzzy.DetTriangleMF(s, sWVar)
	sM, _ := fuzzy.DetTriangleMF(s, sMVar)
	sS, _ := fuzzy.DetTriangleMF(s, sSVar)
	sVs, _ := fuzzy.DetTriangleMF(s, sVsVar)

	//find fuzzy output data (area) : inference method: Max-min method
	set0 := fMin(outputGradesDict[0], sVw)
	set1 := fMin(outputGradesDict[1], sW)
	set2 := fMin(outputGradesDict[2], sM)
	set3 := fMin(outputGradesDict[3], sS)
	set4 := fMin(outputGradesDict[4], sVs)

	//Aggregate all output - Max of Min
	outputArea := fMax(set0, fMax(set1, fMax(set2, fMax(set3, set4))))

	//Defuzzification using CoG
	strength := fuzzy.Defuzzy(s, outputArea)

	return strength
}

//fMax compares two arrays and returns a new array containing the element-wise maxima.
func fMax(a, b []float64) []float64 {
	fuzzy.PanicIfUnequalLength(a, b, "a", "b")
	for i, val := range b {
		if a[i] < val {
			a[i] = val
		}
	}
	return a
}

//fMin compares two arrays and returns a new array containing the element-wise minima.
func fMin(a float64, b []float64) []float64 {
	for i, val := range b {
		if a < val {
			b[i] = a
		}
	}
	return b
}
