package main

import (
	"fmt"
	"github.com/sfreiberg/gotwilio"
	"os"
)

func main() {
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	twilio := gotwilio.NewTwilioClient(accountSid, authToken)

	from := "+16266281388"
	to := "+8617731895913"
	message := "Test for the Message"
	fmt.Println(twilio.SendSMS(from, to, message, "", ""))
}
