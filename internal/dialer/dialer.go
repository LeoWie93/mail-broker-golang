package dialer

import (
	"gopkg.in/gomail.v2"
	"os"
)

var Dialer *gomail.Dialer

func InitDialer() {
	Dialer = gomail.NewDialer(os.Getenv("SMTP_HOST"), 25, os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"))
}
