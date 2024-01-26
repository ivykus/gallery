package main

import (
	"fmt"

	"github.com/go-mail/mail/v2"
)

// Host	sandbox-smtp.mailcatch.app
// Port	25, 1025, 2525
// Username	c8648ea7b08f
// Password	711d678a04e4

const (
	host     = "sandbox-smtp.mailcatch.app"
	port     = 1025
	username = "c8648ea7b08f"
	password = "711d678a04e4"
)

func main() {
	from := "test@ivygallery.com"
	to := "themoonissilent@gmail.com"
	subject := "this is a test email"
	plaintext := "this is a body"
	html := "<h1>Hello</h1><p>This is email</p>"

	msg := mail.NewMessage()
	msg.SetHeader("To", to)
	msg.SetHeader("From", from)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", plaintext)
	msg.AddAlternative("text/html", html)
	// msg.WriteTo(os.Stdout)

	dialer := mail.NewDialer(host, port, username, password)
	err := dialer.DialAndSend(msg)
	if err != nil {
		panic(err)
	}
	fmt.Println("Msg sent")

}
