package metric

import (
	"fmt"
)

// CalculateLength calculates the length (as int) of the given password string.
func CalculateLength(password string) int {
	return len([]rune(password))
}

// GetHintLength provides a hint for a given length
func GetHintLength(length float64, language string) string {
	var message string
	if language == "de" {
		message = fmt.Sprintf("Deine Eingabe hat die Länge %v. ", int(length))
		if length <= 5 {
			message = message + "Das ist sehr schlecht."
		} else if length <= 11 {
			message = message + "Das ist schlecht."
		} else if length <= 17 {
			message = message + "Das ist okay."
		} else if length <= 23 {
			message = message + "Das ist gut."
		} else if length > 23 {
			message = message + "Das ist sehr gut."
		}
	} else {
		message = fmt.Sprintf("Your input has the length %v. ", int(length))
		if length <= 5 {
			message = message + "That's very bad."
		} else if length <= 11 {
			message = message + "That's bad."
		} else if length <= 17 {
			message = message + "It's okay."
		} else if length <= 23 {
			message = message + "That's good."
		} else if length > 23 {
			message = message + "That's very good."
		}
	}
	return message
}
