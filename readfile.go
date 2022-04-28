package main

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

func ReadLicense(path string) string {
	lif("Reading PMTA license file: [%v]", path)
	file, err := os.Open(path)
	if err != nil {
		lef("Failed to open file [%v]", path)
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	regex, err := regexp.Compile(`expires: [0-9]{4}-[0-9]{2}-[0-9]{2}`)

	if err != nil {
		lef("Failed to compile regex [%v]", err)
		return ""
	}

	for scanner.Scan() {
		if regex.MatchString(scanner.Text()) {
			return trimExpire(regex.FindString(scanner.Text()))
		}

		if err := scanner.Err(); err != nil {
			lef("Failed to scan file, %v", err)
			return ""
		}
	}

	return ""
}

func trimExpire(s string) string {
	ldffunc(GetFuncDetails(), "Triming expires prefix from date string [%v]", s)
	return strings.Trim(s, "expires: ")
}
