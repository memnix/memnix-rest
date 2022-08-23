package utils

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"github.com/joho/godotenv"
	gomail "gopkg.in/mail.v2"
	"log"
	"math/big"
	"os"
	"strconv"
)

// GenerateSecretCode generates a secret code
func GenerateSecretCode(length int) string {
	// Generate a random string of letters and digits
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

// GenerateRandomDigit GetRandomNumber returns a random number between min and max
func GenerateRandomDigit(min, max int64) (string, error) {
	// Set max and min value
	randomNumber, err := rand.Int(rand.Reader, big.NewInt(max-min))
	if err != nil {
		return "0", err
	}

	return strconv.FormatInt(randomNumber.Int64()+min, 10), nil
}

// GetSmtpConfig returns a gomail.Dialer and gomail.Message
func getSMTPConfig() (*gomail.Dialer, *gomail.Message) {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	password := os.Getenv("SMTP_PASSWORD")
	host := os.Getenv("SMTP_HOST")
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	from := os.Getenv("SMTP_USER")

	// Settings for SMTP server
	d := gomail.NewDialer(host, port, from, password)
	d.TLSConfig = &tls.Config{
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS12,
		MaxVersion:         0,
		ServerName:         host,
	}
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	return d, m
}

// SendEmail sends an email to the given address
func SendEmail(email, subject, body string) error {
	// Send email

	d, m := getSMTPConfig()

	// Set E-Mail receivers
	m.SetHeader("To", email)

	// Set E-Mail subject
	m.SetHeader("Subject", subject)

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", body)

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}
	return nil
}
