package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Emails struct {
	addresses []string
}

func (e *Emails) GetAddress() []string {
	return e.addresses
}

func (e *Emails) String() string {
	return strings.Trim(fmt.Sprint(e.addresses), "{[]}")
}

func (e *Emails) Set(v string) error {
	if len(e.addresses) > 0 {
		ldffunc(GetFunc(), "Invalid email flag used more than once: %v", "[-e|-email]")
		return fmt.Errorf("Cannot use email flag more than once")
	}
	emails := strings.Split(v, ",")
	for _, item := range emails {
		if IsEmailValid(item) {
			e.addresses = append(e.addresses, item)
		} else {
			ldffunc(GetFunc(), "Invalid email address: %v", item)
			return fmt.Errorf("Invalid email address: %v", item)
		}
	}
	return nil
}

func main() {

	var threshold int
	flag.IntVar(&threshold, "threshold", 5, "the number of days before expiration date")
	flag.IntVar(&threshold, "t", 5, "the number of days before expired date")

	flag.IntVar(&LAST_ALERT, "remainder", 2, "the number of days before expiration date since first alert")
	flag.IntVar(&LAST_ALERT, "r", 2, "the number of days before expired date since first alert")
	var pathtofile string
	flag.StringVar(&pathtofile, "file", "/opt/pmta/license", "path to the license file")
	flag.StringVar(&pathtofile, "f", "/opt/pmta/license", "path to the license file")
	var smtpServer string
	flag.StringVar(&smtpServer, "server", "smtp-server.something.com", "SMTP server name (hostname)")
	flag.StringVar(&smtpServer, "s", "smtp-server.something.com", "SMTP server name (hostname)")
	var smtpPort int
	flag.IntVar(&smtpPort, "port", 25, "port number of smtp server")
	flag.IntVar(&smtpPort, "p", 25, "port number of smtp server")
	var data string
	flag.StringVar(&data, "message", "PMTA License is going to expire in few days", "Email Subject")
	flag.StringVar(&data, "m", "PMTA License is going to expire in few days", "Email Subject")
	var help bool
	flag.BoolVar(&help, "help", false, "display the usage info")
	flag.BoolVar(&help, "h", false, "display the usage info")
	var quiet bool
	flag.BoolVar(&quiet, "quiet", false, "when flag is used, script will only print messages to stdout/log. No emails will be send.")
	flag.BoolVar(&quiet, "q", false, "when flag is used, script will only print messages to stdout/log. No emails will be send.")
	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "when flag is used, script will print debug messages to stdout/log file.")
	flag.BoolVar(&verbose, "v", false, "when flag is used, script will print debug messages to stdout/log file.")
	var chkmode bool
	flag.BoolVar(&chkmode, "chkmode", false, "when flag is used, script can be plugged into sensu / icinga monitoring tool. Handles two types of status: [0 - OK] [2 - ERROR]")
	flag.BoolVar(&chkmode, "c", false, "when flag is used, script can be plugged into sensu / icinga monitoring tool. Handles two types of status: [0 - OK] [2 - ERROR]")
	var emails Emails
	flag.Var(&emails, "email", "email addresses to which we are going to send notifications. Comma separated values.")
	flag.Var(&emails, "e", "email addresses to which we are going to send notifications. Comma separated values.")

	flag.Parse()

	InitLoggerStdoutOnly(chkmode)
	SetDebug(verbose)
	InitLogger()

	ldln("=======")
	ldln("Passed parameters:")
	ldf("[-h|-help] <%v>", help)
	ldf("[-e|-email addresses] <%v>", emails.String())
	ldf("[-t|-threshold num] <%v>", threshold)
	ldf("[-r|-remainder num] <%v>", LAST_ALERT)
	ldf("[-f|-file pathToFile] <%v>", pathtofile)
	ldf("[-s|-server smtp host] <%v>", smtpServer)
	ldf("[-p|-port port] <%v>", smtpPort)
	ldf("[-m|-message text] <%v>", data)
	ldf("[-v|-verbose] <%v>", verbose)
	ldln("=======")

	if threshold < LAST_ALERT {
		lef("Remainder cannot be greater than threshold. [%v > %v] == [%v]", threshold, LAST_ALERT, (threshold > LAST_ALERT))
		Usage(threshold, LAST_ALERT)
	}

	if help {
		Usage(threshold, LAST_ALERT)
	}

	if quiet {
		MainStdOutput(threshold, pathtofile)
	} else if chkmode {
		MainCheckMode(threshold, pathtofile)
	} else {
		Main(threshold, smtpPort, pathtofile, smtpServer, data, emails.addresses...)
	}

}

func Usage(threshold, LAST_ALERT int) {
	filename := filepath.Base(os.Args[0])
	cmd := `

Usage: 
	%v [-h|-help] [-e|-email addresses] [-t|-threshold num] [-r|-remainder num] [-f|-file pathToFile] [-s|-server smtp host] [-p|-port port] [-m|-message text] [-q|-quiet] [-v|-verbose]

Description:
	`
	pf(cmd+"\n", filename)
	flag.PrintDefaults()
	usage := `
The script will check the license for PMTA Server, and if it is below the threshold
Mail will be send for the notification
There are going to be 2 alerts before the expire date, and continuous alerts after expire date.
	- 1st alert will be triggered at X - [%[1]v] days equal threshold
	- 2nd alert will be triggered at X - Y (%[1]v - %[3]v = %[2]v) days before expire date

**Note: 
	X is the threshold number
	Y is the remainder number
	where Y < X
	`
	pf(usage+"\n", threshold, (threshold - LAST_ALERT), LAST_ALERT)
	os.Exit(0)
}
