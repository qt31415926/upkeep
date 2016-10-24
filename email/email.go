package email

import (
	"fmt"
	"net/smtp"
)

func SmtpSend(from, pass, to, url string) {
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + "Something is wrong with: " + url + "\n\n" +
		"The site may be down, or displaying the wrong content"
	err := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", from, pass, "smtp.gmail.com"), from, []string{to}, []byte(msg))
	if err != nil {
		fmt.Printf("Error, %s", err)
	}

	return
}
