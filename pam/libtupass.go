package main

import (
	"C"
	"log"

	"github.com/tupass/tupass-backend/api"
)

// default value for password list, can be overridden in build script
var pwlist = "10-million-password-list-top-50000.txt"

//CalculateStrength calculates the total strength for given password
//export CalculateStrength
func CalculateStrength(password *C.char) (s float64) {
	// defer panic error that occurs when setup.go cannot read password file
	defer func() {
		if r := recover(); r != nil {
			log.Println("Unable to load TUPass password list: ", r)
			s = 0
		}
	}()

	// read password file, calculate and return result
	api.SetupPasswordByFile(pwlist)
	_, _, _, s, _, _, _, _ = api.CalculateMetrics(C.GoString(password))
	return
}

func main() {
}
