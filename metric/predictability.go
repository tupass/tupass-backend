package metric

import (
	"fmt"
	"unicode"
)

// PasswordList is an array of []rune(s) to store the password list in the programs heap.
var PasswordList [][]rune

// empty is dummy variable (for shortage) for type empty struct.
var empty struct{}

// leetspeakMap is HashMap, containing HashMaps used to translates relevant lower-/uppercase letters to special ascii characters looking similar.
// Using a normal letter as key, contained value is also a HashMap containing similar special chars and numbers, looking similar (leet version) to the letter.
// The second HashMap is just used as a HashSet for quick occur checking.
var leetspeakMap = map[rune]map[rune]struct{}{
	'a': {'4': empty, '@': empty},
	'b': {'8': empty},
	'c': {'(': empty, '{': empty, '[': empty, '<': empty},
	'e': {'3': empty},
	'g': {'6': empty, '9': empty},
	'i': {'1': empty, '!': empty, '|': empty},
	'l': {'1': empty, '|': empty, '7': empty},
	'o': {'0': empty},
	's': {'$': empty, '5': empty},
	't': {'+': empty, '7': empty},
	'x': {'%': empty},
	'z': {'2': empty},
	'A': {'4': empty, '@': empty},
	'B': {'8': empty},
	'C': {'(': empty, '{': empty, '[': empty, '<': empty},
	'E': {'3': empty},
	'G': {'6': empty, '9': empty},
	'I': {'1': empty, '!': empty, '|': empty},
	'L': {'1': empty, '|': empty, '7': empty},
	'O': {'0': empty},
	'S': {'$': empty, '5': empty},
	'T': {'+': empty, '7': empty},
	'X': {'%': empty},
	'Z': {'2': empty}}

// leetCheck compares char1 and char2 on leet-similarity by looking up the presence of char2 in char1's HashMap (see leetspeakMap) and vice versa.
// Only returns true if char1 / char2 occurs in the HashMap of char2 / char1.
func leetCheck(char1, char2 rune) bool {
	if val, ok := leetspeakMap[char1]; ok {
		// for char1 exist leet chars
		if _, ok := val[char2]; ok {
			// char2 is a leet char of char1
			return true
		}
	}
	if val, ok := leetspeakMap[char2]; ok {
		if _, ok := val[char1]; ok {
			return true
		}
	}
	return false
}

// min returns the minimum value of a, b and c.
func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
	} else {
		if b < c {
			return b
		}
	}
	return c
}

// calculateDistance calculates and returns the levenshtein distance between strings a and b.
// aLength is required so that the the length of string a is not needed to be recalculated.
// aColumn is required for an efficient memory usage of the algorithm by providing a pre-constructed int vector of size aLenght.
// This implemention is optimized to use O(a) space.
// Levenshtein distance between two words is the minimum number of single-character edits (insertions, deletions or substitutions)
// required to change one word into the other (see http://en.wikipedia.org/wiki/Levenshtein_distance).
func calculateDistance(a, b []rune, aLength int, aColumn []int) int {
	var cost, lastDiagonalValue, oldDiagonalValue int

	// calculate length of b
	bLength := len(b)

	// initialize aColumn ascending (first columns values)
	for y := 1; y <= aLength; y++ {
		aColumn[y] = y
	}

	for x := 1; x <= bLength; x++ { // outer loop moving column to the right
		// updates to current column's first row entry (simply the number of the column)
		aColumn[0] = x

		lastDiagonalValue = x - 1

		for y := 1; y <= aLength; y++ { // inner loop updating current column
			oldDiagonalValue = aColumn[y]

			// now comparing characters at same index in both strings
			char1, char2 := a[y-1], b[x-1]

			if char1 == char2 {
				// same char -> change distance is 0
				cost = 0
			} else if unicode.ToLower(char1) == unicode.ToLower(char2) || leetCheck(char1, char2) {
				// char just up/down shifted or similar to leet -> change distance is 1
				cost = 1
			} else {
				// char is completely different -> change distance is 2
				cost = 2
			}

			// current cells value gets assigned the minimum of:
			//  previous horizontal left cell's value  +1
			//  previous vertical above cell's value +1
			//  previous diagonal up left cell's value + cost (0: similar, 1: upper/lower change OR leet translation, 2: else)
			aColumn[y] = min(aColumn[y]+1, aColumn[y-1]+1, lastDiagonalValue+cost)
			lastDiagonalValue = oldDiagonalValue
		}
	}

	// return value of lower right corner of (virtual) matrix, containing the levenshtein distance
	return aColumn[aLength]
}

//CalculatePredictability calculates the predictability of the basePassword with the given passwordList
func CalculatePredictability(basePasswordString string) (float64, string) {

	// translate string to rune array
	basePassword := []rune(basePasswordString)
	// calculate length of basePassword
	basePasswordLength := len(basePassword)
	// create column for optimized levenshtein calculation
	basePasswordColumn := make([]int, basePasswordLength+1)

	// initialize similarity with zero
	greatestSimilarity := float64(0)
	mostSimilarPassword := ""

	// iterate over every password in passwordList to calc distance and the resulting similarity
	for _, currentPassword := range PasswordList {

		distance := calculateDistance(basePassword, currentPassword, basePasswordLength, basePasswordColumn)
		lengthSum := float64(basePasswordLength + len(currentPassword))

		// see slide 23 of theory presentation
		currentSimilarity := 1 - float64(distance)/lengthSum

		// only cosider greatest similarity
		if currentSimilarity > greatestSimilarity {
			// update similarity
			greatestSimilarity = currentSimilarity
			mostSimilarPassword = string(currentPassword)
		}
	}

	//  P = max(Similarity to username, Similarity to common list)
	// currently no username -> predictability = similarity to common list
	return greatestSimilarity * 100, mostSimilarPassword
}

// GetHintPredictability provides the most similar password as a hint if its predictability is higher than 50
func GetHintPredictability(mostSimilarPassword string, score float64, language string) string {
	if language == "de" {
		if score > 80 {
			return fmt.Sprintf("Dein Passwort ist sehr ähnlich zu '%s' in unserer Passwortliste.", mostSimilarPassword)
		} else if score > 60 {
			return fmt.Sprintf("Dein Passwort ist ähnlich zu '%s' in unserer Passwortliste.", mostSimilarPassword)
		} else if score > 40 {
			return fmt.Sprintf("Die Ähnlichkeit deines Passwortes zu den Passwörtern in unserer Liste ist gering. Gut!")
		} else if score > 20 {
			return fmt.Sprintf("Die Ähnlichkeit deines Passwortes zu den Passwörtern in unserer Liste ist sehr gering. Sehr gut!")
		}
		return "Wir haben kein ähnliches Passwort in unserer Liste gefunden. Gute Arbeit!"
	}
	if score > 80 {
		return fmt.Sprintf("Your password is very similar to '%s' in our password list.", mostSimilarPassword)
	} else if score > 60 {
		return fmt.Sprintf("Your password is similar to '%s' in our password list.", mostSimilarPassword)
	} else if score > 40 {
		return fmt.Sprintf("The similarity of your password to the passwords in our list is low. Good!")
	} else if score > 20 {
		return fmt.Sprintf("The similarity of your password to the passwords in our list is very low. Great!")
	}
	return "No similar password was found in our list. Good job!"
}
