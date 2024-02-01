	package otp

	import (
		"crypto/rand"
		"encoding/base64"
		"fmt"
		"net/smtp"
		"os"

		mail "github.com/jordan-wright/email"
	)

	func GenerateOtp() (string, error) {
		randomBytes := make([]byte, 3)
		_, err := rand.Read(randomBytes)
		if err != nil {
			fmt.Print(err)
			return "", err
		}
		return base64.URLEncoding.EncodeToString(randomBytes),nil
	}

	func SendOTP(email []string,otp string) error {
		e := mail.NewEmail()
		from := os.Getenv("sender")
		appPassword := os.Getenv("appPassword")
		name := os.Getenv("name")
		e.From = fmt.Sprintf("%s <%s>",name,from)
		msg := "Subject: Verify your account\n Your code is " + otp
		auth := smtp.PlainAuth(
			"",
			from,
			appPassword,
			"smtp.gmail.com",
		)

		err := smtp.SendMail(
			"smtp.gmail.com:587",
			auth,
			from,
			email,
			[]byte(msg),
		)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("Sending OTP %s to %s\n",otp,email)
		return nil
	}