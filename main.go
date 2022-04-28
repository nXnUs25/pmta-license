package main

import (
	"fmt"
	"os"
	"time"
)

const (
	LOCK_PATH      = "/tmp/pmta-license-check.lock"
	LAST_LOCK_PATH = "/tmp/pmta_license_last_alert.lock"
)

var LAST_ALERT = 2

func Main(threshold, port int, path, server, data string, emails ...string) {
	date := ReadLicense(path)
	days := ExpireInDays(date)

	template := `
	Please update your license for PMTA [%v] Server [%v]
	Because will expire in %v days.
	
	Regards
	SRE
	`
	msg := fmt.Sprintf(template, date, GetHostname(), days)
	if SpamThem(days) {
		lif("SpamThem because no %v days left", days)
		if CheckIfLockExists(LOCK_PATH) {
			RemoveLocks(LOCK_PATH)
		}
		if CheckIfLockExists(LAST_LOCK_PATH) {
			RemoveLocks(LAST_LOCK_PATH)
		}
	}

	if days > threshold {
		lif("Clear Locks because %v days are above %v threshold, License Updated!", days, threshold)
		if CheckIfLockExists(LOCK_PATH) {
			RemoveLocks(LOCK_PATH)
		}
		if CheckIfLockExists(LAST_LOCK_PATH) {
			RemoveLocks(LAST_LOCK_PATH)
		}
	}

	if Alert(days, threshold) && (threshold-LAST_ALERT) < days {
		lif("Alert because %v days left", days)
		if !CheckIfLockExists(LOCK_PATH) {
			LockAlerts(LOCK_PATH)
			SendEmail(port, server, data, msg, emails...)
			lif("Lock %v set and email to %v sent because %v days left", LOCK_PATH, emails, days)
		} else {
			LockStillThere(LOCK_PATH)
			lif("Lock set %v, %v days left low", LOCK_PATH, days)
		}
	} else if LastAlert(days, threshold) {
		lif("LastAlert because %v days left", days)
		if !CheckIfLockExists(LAST_LOCK_PATH) {
			LockAlerts(LAST_LOCK_PATH)
			SendEmail(port, server, data, msg, emails...)
			lif("Last Lock %v set and Last email to %v sent because %v days left", LAST_LOCK_PATH, emails, days)
		} else {
			LockStillThere(LAST_LOCK_PATH)
			lif("Last Lock set %v, %v days left low", LAST_LOCK_PATH, days)
		}
	} else {
		lif("All seems to be good, %v days left to expire date, which is higher than threshold %v", days, threshold)
	}
}

func MainStdOutput(threshold int, path string) {
	liln("Running check in quiet mode, Alerts will not be send, just printed to stdout")
	date := ReadLicense(path)
	days := ExpireInDays(date)

	template := `
	Please update your license for PMTA [%v] Server [%v]
	Because will expire in %v days.
	
	Regards
	SRE
	`
	msg := fmt.Sprintf(template, date, GetHostname(), days)
	if Alert(days, threshold) {
		lif("Alert because %v days left, which is lower than threshold %v", days, threshold)
		lif("Alert message:\n %v", msg)
		liln("Quiet mode check finished.")
		os.Exit(2)
	} else {
		lif("No Alerts, %v days left to expire date [%v], which is higher than threshold %v", days, date, threshold)
		liln("Quiet mode check finished.")
		os.Exit(0)
	}
}

func MainCheckMode(threshold int, path string) {
	date := ReadLicense(path)
	days := ExpireInDays(date)

	template := `License for PMTA will expire after [%v] on Server [%v] days left till expire date is %v days.`
	msg := fmt.Sprintf(template, date, GetHostname(), days)
	SetStdout()
	if Alert(days, threshold) {
		lef("[CRITICAL] - %v", msg)
		os.Exit(2)
	} else {
		lif("[OK] - %v days left to expire date [%v], which is higher than threshold %v", days, date, threshold)
		os.Exit(0)
	}
}

func Alert(days, threshold int) bool {
	return days <= threshold
}

func GetHostname() string {
	node_name, err := os.Hostname()
	if err != nil {
		return ""
	}
	return node_name
}

func LockAlerts(path string) {
	lif("LockAlerts %v", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	checkError("OpenFile", err)
	defer f.Close()
	t := time.Now()
	content := "Lock created at " + t.String() + "\n"
	_, err = f.Write([]byte(content))
	checkError("OpenFile", err)
}

func LockStillThere(path string) {
	lif("LockStillThere %v ", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	checkError("LocksStilThere", err)
	defer f.Close()
	t := time.Now()
	content := "Lock existed and accessed at " + t.String() + "\n"
	_, err = f.Write([]byte(content))
	checkError("Write", err)
}

func CheckIfLockExists(path string) bool {
	lif("CheckIfLockExists %v ", path)
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func RemoveLocks(path string) {
	lif("RemoveLocks %v ", path)
	if err := os.Remove(path); err != nil {
		checkError("Remove", err)
	}
}

func SpamThem(number int) bool {
	return number <= 0
}

func LastAlert(number int, threshold int) bool {
	return Alert(number, threshold-LAST_ALERT)
}
