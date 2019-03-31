package api

import (
	"log"
	"math"

	"github.com/tupass/tupass-backend/fes"
	"github.com/tupass/tupass-backend/fuzzy"
	"github.com/tupass/tupass-backend/metric"
)

// MetricResult is a struct representing a metric result provided to a client
type MetricResult struct {
	Score   int    `json:"score"`
	Message string `json:"message"`
	Hint    string `json:"hint"`
}

// Result is a struct representing all model calculation results provided to a client
type Result struct {
	Length         MetricResult `json:"length"`
	Complexity     MetricResult `json:"complexity"`
	Predictability MetricResult `json:"predictability"`
	Strength       MetricResult `json:"strength"`
}

// CalculateMetrics calculates the results length, complexity, predictability, total strength, corresponding membership grades and the mostSimilarPassword for a given password string
func CalculateMetrics(password string) (length, complexity, predictability, strength float64, LList, CList, PList []float64, mostSimilarPassword string) {
	// calculate the main metrics
	length = float64(metric.CalculateLength(password))
	complexity = metric.CalculateComplexity(password)
	predictability, mostSimilarPassword = metric.CalculatePredictability(password)

	// calculate memberships of metric values
	LList = fuzzy.CalculateMembershipGradesForLength(length)
	CList = fuzzy.CalculateMembershipGradesForComplexity(complexity)
	PList = fuzzy.CalculateMembershipGradesForPredictability(predictability)

	// calculate overall strength
	strength = fes.GetStrengthByMembershipGrades(LList, CList, PList)
	return
}

// CalculateResult calculates the results and provides a Result struct representation of the length, complexity, predictability and total strength for a given password string
func CalculateResult(password string, language string) Result {
	length, complexity, predictability, strength, LList, CList, PList, mostSimilarPassword := CalculateMetrics(password)

	return Result{
		Length:         getLengthResult(length, LList, language),
		Complexity:     getComplexResult(complexity, CList, password, language),
		Predictability: getPredictabilityResult(predictability, PList, mostSimilarPassword, language),
		Strength:       getStrengthResult(strength, language)}
}

// getStrengthScore provides a MetricResult struct representation of given total strength
func getStrengthResult(strength float64, language string) MetricResult {
	strengthScore := int(math.Round(strength))
	strengthMsg := strengthResultToText(strength, language)
	return MetricResult{
		Score:   strengthScore,
		Message: strengthMsg,
		Hint:    ""}
}

// strengthResultToText turns a strength level in float64 (percent, e.g. 20.43523) to the corresponding set name.
func strengthResultToText(level float64, language string) string {
	level = math.Round(level)
	var textResult string

	if language == "de" {
		if 0 <= level && level <= 20 {
			textResult = "sehr schwach"
		} else if 20 < level && level <= 40 {
			textResult = "schwach"
		} else if 40 < level && level <= 60 {
			textResult = "mittelmäßig"
		} else if 60 < level && level <= 80 {
			textResult = "stark"
		} else if 80 < level && level <= 100 {
			textResult = "sehr stark"
		} else {
			textResult = "no such level"
		}
	} else {
		if 0 <= level && level <= 20 {
			textResult = "very weak"
		} else if 20 < level && level <= 40 {
			textResult = "weak"
		} else if 40 < level && level <= 60 {
			textResult = "medium"
		} else if 60 < level && level <= 80 {
			textResult = "strong"
		} else if 80 < level && level <= 100 {
			textResult = "very strong"
		} else {
			textResult = "no such level"
		}
	}

	return textResult
}

// getLengthScore provides a MetricResult struct representation of the given length and length membership grades
func getLengthResult(length float64, LList []float64, language string) MetricResult {
	// A password is very long (->100%) when it has 26 or more characters. => maxVal=26
	var linguisticVars []string
	if language == "de" {
		linguisticVars = []string{"sehr kurz", "kurz", "mittelmäßig", "lang", "sehr lang"}
	} else {
		linguisticVars = []string{"very short", "short", "medium", "long", "very long"}
	}

	hint := metric.GetHintLength(length, language)
	return generateMetricResult(length, 26, LList, linguisticVars, hint, language)
}

// getComplexScore provides a MetricResult struct representation of the given complexity and complecity membership grades
func getComplexResult(complexity float64, CList []float64, password string, language string) MetricResult {
	// A password is very complex (->100%) when it has 26 or more characters. => maxVal=677
	var linguisticVars []string
	if language == "de" {
		linguisticVars = []string{"sehr einfach", "einfach", "mittelmäßig", "komplex", "sehr komplex"}
	} else {
		linguisticVars = []string{"very simple", "simple", "medium", "complex", "very complex"}
	}

	hint := metric.GetHintComplexity(password, complexity, language)
	return generateMetricResult(complexity, 677, CList, linguisticVars, hint, language)
}

// getPredictabilityScore provides a MetricResult struct representation of the given predictability and predictability membership grades
func getPredictabilityResult(predictability float64, PList []float64, mostSimilarPassword string, language string) MetricResult {
	var linguisticVars []string
	if language == "de" {
		linguisticVars = []string{"schwer vorherzusagen", "mittelmäßig", "einfach vorherzusagen"}
	} else {
		linguisticVars = []string{"hard to predict", "medium", "easy to predict"}
	}

	hint := metric.GetHintPredictability(mostSimilarPassword, predictability, language)
	return generateMetricResult(predictability, 100, PList, linguisticVars, hint, language)
}

// generateMetricResult returns a MetricResult struct representation of a metric,
// given its float64 value, maximum value, membership grades and lingustic variables.
// grades and linguisticVars must have the same dimensions.
// It normalizes and rounds the metric value in order to provide the metrics score as percentage (int)
// and generates a message to representing the membership grades using lingusticVariables.
func generateMetricResult(value float64, maxValue float64, grades []float64, linguisticVars []string, hint string, language string) MetricResult {
	// normalize and round metric value
	var score int
	if value > maxValue {
		score = 100
	} else {
		score = int(math.Round(value / maxValue * 100))
	}

	// generate message using membership grades
	if len(grades) != len(linguisticVars) {
		log.Panicf("Expected grades and linguisticVars to have the same length, but found lengths of %d and %d instead.", len(grades), len(linguisticVars))
	}
	maxGradeIdx := getMaxGradeIndex(grades)
	checkSecondMax := checkSecondMaxGrade(grades, maxGradeIdx)
	var msg = linguisticVars[maxGradeIdx]
	if checkSecondMax != -1 { //if we have two maxima adapt message
		if language == "en" {
			msg += " to "
		} else {
			msg += " bis "
		}
		msg += linguisticVars[checkSecondMax]
	}

	return MetricResult{Score: score, Message: msg, Hint: hint}
}

// getMaxGradeIndex returns the index of the maximum value in a given array.
// It is used to determine the best membership grade.
func getMaxGradeIndex(list []float64) (idx int) {
	var max = list[0]
	idx = 0
	for i, grade := range list {
		if grade > max {
			max = grade
			idx = i
		}
	}
	return
}

// checkSecondMaxGrade checks whether or not there is a second element in a given array that has the same value as the given index,
// It returns -1 if no such element is found, otherwise the index is returned.
// It is used to check if two memberships apply at the same time.
func checkSecondMaxGrade(list []float64, idx int) int {
	max := list[idx]
	for i, grade := range list {
		if grade == max && i != idx {
			return i
		}
	}
	return -1
}
