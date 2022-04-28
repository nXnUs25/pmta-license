package main

import (
	"bytes"
	"fmt"
	"net/smtp"
	"regexp"
)

func checkError(msg string, err error) {
	if err != nil {
		lef("%v : %v", msg, err)
	}
}

func SendEmail(smtpPort int, smtpServer, subject, msg string, to ...string) {
	host := GetHostname()

	from := "pmta_license@" + host

	smtpHost := fmt.Sprintf("%v:%v", smtpServer, smtpPort)
	lif("Opening connection to SMTP Server %v", smtpHost)
	hfrom := "From: " + from
	done := make(chan bool, len(to))
	for _, t := range to {

		go func(t string) {
			conn, err := smtp.Dial(smtpHost)
			checkError("Dialing", err)
			defer conn.Close()
			err = conn.Mail(from)
			checkError("Mail", err)
			err = conn.Rcpt(t)
			checkError("Rcpt", err)
			wc, err := conn.Data()
			checkError("Data", err)
			defer wc.Close()
			hto := "To: " + t
			hsub := "Subject: " + subject
			body := hfrom + "\n" + hto + "\n" + hsub + "\n" + msg
			buffer := bytes.NewBufferString(body)
			if _, err := buffer.WriteTo(wc); err != nil {
				leln(err)
			}
			lif("Notification sent to [%v]", t)
			done <- true
		}(t)
		<-done
	}
	lif("Closing connection to SMTP Server %v", smtpHost)

}

func IsEmailValid(e string) bool {
	ldffunc(GetFuncDetails(), "email address: %v", e)
	emailRegexPattern := "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	emailRegex := regexp.MustCompile(emailRegexPattern)
	return emailRegex.MatchString(e)
}
