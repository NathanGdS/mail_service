package main

import (
	"log"
	"net/smtp"
	"os"

	gomail "gopkg.in/gomail.v2"

	"github.com/joho/godotenv"
)

func getEnvVariable(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func sendMailWithSmtp(subject string, body string, to []string) {
	gmailAddress := getEnvVariable("EMAIL_ADDRESS")
	host := getEnvVariable("EMAIL_HOST")
	appPassword := getEnvVariable("EMAIL_APP_PASSWORD")
	from := getEnvVariable("FROM_EMAIL")

	messageBody := []byte("Subject: " + subject + "\r\n\r\n" + body)

	auth := smtp.PlainAuth(
		"",
		from,
		appPassword,
		host,
	)

	err := smtp.SendMail(
		gmailAddress,
		auth,
		from,        // From
		to,          // To
		messageBody, // Body Message
	)

	if err != nil {
		log.Fatalln(err)
	}
}

type list_type []struct {
	Address string
}

func sendMailWithGoMail(subject string, body string, to list_type) {
	host := getEnvVariable("EMAIL_HOST")
	dialer := gomail.NewDialer(host, 587, getEnvVariable("FROM_EMAIL"), getEnvVariable("EMAIL_APP_PASSWORD"))

	s, err := dialer.Dial()
	if err != nil {
		panic(err)
	}

	mailer := gomail.NewMessage()
	for _, r := range to {
		mailer.SetHeader("From", getEnvVariable("FROM_EMAIL"))
		mailer.SetHeader("To", r.Address)
		mailer.SetHeader("Subject", subject)
		mailer.SetBody("text/html", body)

		if err := gomail.Send(s, mailer); err != nil {
			log.Printf("Could not send email to %q: %v", r.Address, err)
		}
		mailer.Reset()
		toLog := "Email sent to \"" + r.Address + "\" usign GoMail"
		log.Printf(toLog)
	}
}

func main() {
	log.Println("Starting mailer...")

	body := string("This is a test email")
	// targetEmails := []string{""} // target emails
	// sendMailWithSmtp(string("Sended with Smtp"), body, targetEmails)
	sendMailWithGoMail(string("Sended with GoMail"), body,
		list_type{
			{Address: ""}, //list of target emails
		})

	log.Println("Mailer finished...")
}
