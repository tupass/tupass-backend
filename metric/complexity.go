package metric

import (
	"log"
	"unicode"
)

// numberOfChars calculates the number of characters
func numberOfChars(pw string) (int, int, int, int) {
	var low = 0
	var up = 0
	var d = 0
	var special = 0

	for _, r := range pw {
		switch true {
		case r < 32 || r == 127:
			{
				log.Println("Can not handle ASCII characters < 32")
				return -1, -1, -1, -1
			}
		case unicode.IsLower(r):
			low++
		case unicode.IsUpper(r):
			up++
		case unicode.IsNumber(r):
			d++
		default:
			special++
		}
	}

	return low, up, d, special
}

// CalculateComplexity calculates the complexity for a given password string.
// It returns the complexity as float64.
func CalculateComplexity(pw string) float64 {
	low, up, d, special := numberOfChars(pw)

	complexity := low*26 + up*26 + d*10 + special*33
	complexityFloat := ComplexityReward(complexity, low, up, d, special)

	return complexityFloat
}

// ComplexityReward calculates the rewards for given complexity c and numbers of runes n, m, p and q of each type.
// It returns c mulitpilied with  the number of types that unequal 0.
func ComplexityReward(c int, n int, m int, p int, q int) float64 {
	temp := [4]int{n, m, p, q}
	count := .0
	for i := range temp {
		if temp[i] != 0 {
			count++
		}
	}

	return float64(c) * 0.25 * count
}

// GetHintComplexity provides a hint for the metric complexity based on the password and its complexity
func GetHintComplexity(pw string, complexity float64, language string) string {
	if language == "de" {
		if complexity > 543.0 { // exactly between complex and very complex, no hint necessary
			return "Die Komplexität deines Passwortes ist sehr gut."
		} else if complexity > 375.0 { // exactly between medium and complex, no hint necessary
			return "Die Komplexität deines Passwortes ist gut."
		}
	} else {
		if complexity > 543.0 { // exactly between complex and very complex, no hint necessary
			return "Your password's complexity is very good."
		} else if complexity > 375.0 { // exactly between medium and complex, no hint necessary
			return "Your password's complexity is good."
		}
	}

	low, up, d, special := numberOfChars(pw)
	total := low + up + d + special

	// password has length 0
	if total <= 0 {
		if language == "de" {
			return "Das Passwort enthält keine Zeichen."
		}
		return "The password has no characters."
	}

	lowPerc := float64(low) / float64(total)
	upPerc := float64(up) / float64(total)
	dPerc := float64(d) / float64(total)
	specialPerc := float64(special) / float64(total)

	makeHint := [4]bool{false, false, false, false}
	hintTotal := 0

	// Decide which characters should be more present
	if lowPerc == 0 || lowPerc < 0.125 {
		makeHint[0] = true
		hintTotal++
	}
	if upPerc == 0 || upPerc < 0.125 {
		makeHint[1] = true
		hintTotal++
	}
	if dPerc == 0 || dPerc < 0.08 {
		makeHint[2] = true
		hintTotal++
	}
	if specialPerc == 0 || specialPerc < 0.125 {
		makeHint[3] = true
		hintTotal++
	}

	// no hint necessary
	if language == "de" {
		if hintTotal <= 0 && complexity > 375 {
			return "Deine Komplexität ist gut."
		} else if hintTotal <= 0 && complexity > 341 {
			return "Deine Komplexität ist gut, aber kann verbessert werden. Vielleicht solltest du mehr Zeichen verwenden."
		} else if hintTotal <= 0 && complexity < 173 {
			return "Du nutzt jeden Zeichensatz, aber nicht genügend Zeichen von jedem."
		} else if hintTotal <= 0 && complexity <= 341 {
			return "Du nutzt jeden Zeichensatz, aber generell könntest Du noch mehr Zeichen verwenden."
		}
	} else {
		if hintTotal <= 0 && complexity > 375 {
			return "Your complexity is good."
		} else if hintTotal <= 0 && complexity > 341 {
			return "Your complexity is ok, but can be improved. Maybe you should use more characters."
		} else if hintTotal <= 0 && complexity < 173 {
			return "You use all sets of characters, but not enough from each."
		} else if hintTotal <= 0 && complexity <= 341 {
			return "You use all sets of characters, but in general you could use even more."
		}
	}

	var message string
	if language == "de" {
		message = "Dein Passwort enthält sehr wenige "
		message = message + getMissingCharactersDE(makeHint, hintTotal)
	} else {
		message = "Your password has very few "
		message = message + getMissingCharactersEN(makeHint, hintTotal)
	}

	return message
}

func getMissingCharactersEN(makeHint [4]bool, hintTotal int) string {
	hintsCovered := 0
	message := ""
	// composition of string
	for i := 0; i < 4; i++ {
		if makeHint[i] == true && hintsCovered < hintTotal-2 {
			hintsCovered++
			if i == 0 {
				message = message + "lowercase letters, "
			} else if i == 1 {
				message = message + "uppercase letters, "
			}
		} else if makeHint[i] == true && hintsCovered < hintTotal-1 {
			hintsCovered++
			if i == 0 {
				message = message + "lowercase letters and "
			} else if i == 1 {
				message = message + "uppercase letters and "
			} else if i == 2 {
				message = message + "digits and "
			}
		} else if makeHint[i] == true && hintsCovered < hintTotal {
			hintsCovered++
			if i == 0 {
				message = message + "lowercase letters."
			} else if i == 1 {
				message = message + "uppercase letters."
			} else if i == 2 {
				message = message + "digits."
			} else if i == 3 {
				message = message + "special characters."
			}
		}
	}
	return message
}

func getMissingCharactersDE(makeHint [4]bool, hintTotal int) string {
	hintsCovered := 0
	message := ""
	// composition of string
	for i := 0; i < 4; i++ {
		if makeHint[i] == true && hintsCovered < hintTotal-2 {
			hintsCovered++
			if i == 0 {
				message = message + "Kleinbuchstaben, "
			} else if i == 1 {
				message = message + "Großbuchstaben, "
			}
		} else if makeHint[i] == true && hintsCovered < hintTotal-1 {
			hintsCovered++
			if i == 0 {
				message = message + "Kleinbuchstaben und "
			} else if i == 1 {
				message = message + "Großbuchstaben und "
			} else if i == 2 {
				message = message + "Ziffern und "
			}
		} else if makeHint[i] == true && hintsCovered < hintTotal {
			hintsCovered++
			if i == 0 {
				message = message + "Kleinbuchstaben."
			} else if i == 1 {
				message = message + "Großbuchstaben."
			} else if i == 2 {
				message = message + "Ziffern."
			} else if i == 3 {
				message = message + "Sonderzeichen."
			}
		}
	}
	return message
}
