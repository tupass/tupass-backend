package api

import (
	"bufio"
	"log"
	"path"

	"github.com/tupass/tupass-backend/metric"

	rice "github.com/GeertJohan/go.rice"
)

// SetupPasswordList calls SetupPasswordByFile with default password list
func SetupPasswordList() {
	filepath := "10-million-password-list-top-50000.txt"
	SetupPasswordByFile(filepath)
}

// SetupPasswordByFile reads given password list (in folder passwords) to memory and appends metric.PasswordList for usage in predictability later on
func SetupPasswordByFile(pwlist string) {
	// first read password list from filepath to memory
	box, err := rice.FindBox("../passwords")
	if err != nil {
		log.Panicf("Could not find directory containing password lists %s\n", err)
	}

	filename := path.Base(pwlist)
	file, err := box.Open(filename)
	if err != nil {
		log.Panicf("Could not open passwordlist %s\n", err)
	}

	// iterate over all lines in file and append them to metric.PasswordList
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := []rune(scanner.Text())
		metric.PasswordList = append(metric.PasswordList, line)
	}

	err = scanner.Err()
	if err != nil {
		log.Panicf("Error while scanning passwordlist %s\n", err)
	}

	err = file.Close()
	if err != nil {
		log.Panicf("Could not close passwordlist %s\n", err)
	}

	log.Printf("Reading password list done.\n")
}
