package main

import (
	"math"
	"time"
)

func ExpireInDays(date string) int {
	lif("Calculating days to expire date [%v]", date)
	ldffunc(GetFuncDetails(), "Calculating days to expire date [%v]", date)
	end_time, err := time.Parse("2006-01-02", date)
	ldf("Date is:[%v] Year:[%v] Month:[%v] Day:[%v]", date, end_time.Year(), end_time.Month(), end_time.Day())
	if err != nil {
		lef("Cannot parse expire date [%v] with - %v", end_time, err)
	}
	today := time.Now()
	days := math.RoundToEven(end_time.Sub(today).Hours() / 24)
	lif("License will expire after %v days", days)
	return int(days)
}
