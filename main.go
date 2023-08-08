package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"net/smtp"
	"strconv"
	"github.com/joho/godotenv"
	// "log"
	// "strings"
)

type Config struct {
	Email string
	Name string
	Password string
	Host string
	Port int
}

func getEnv(key string) string {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading.env file")
    }
    return os.Getenv(key)
}

func sendMail(conf Config, to []string, cc []string, subject, message string) error {
	body := "From: " + conf.Name + "\n" +
		    "To: " + strings.Join(to, ",") + "\n" +
		    "Cc: " + strings.Join(cc, ",") + "\n" +
		    "Subject: " + subject + "\n\n" +
		    message

	auth := smtp.PlainAuth("", conf.Email, conf.Password, conf.Host)
	smtpAddr := fmt.Sprintf("%s:%d", conf.Host, conf.Port)

	err := smtp.SendMail(smtpAddr, auth, conf.Email, append(to, cc...), []byte(body))
	if err != nil {
		return err
	}

	return nil
}

func main() {
	var conf Config

	conf.Email = getEnv("CONFIG_AUTH_EMAIL")
	conf.Name = getEnv("CONFIG_SENDER_NAME")
	conf.Password = getEnv("CONFIG_AUTH_PASSWORD")
	conf.Host = getEnv("CONFIG_SMTP_HOST")
	conf.Port, _ = strconv.Atoi(getEnv("CONFIG_SMTP_PORT"))

	to := []string{"ptadityamahendrap@gmail.com"}
	cc := []string{"whatupbiatch69@gmail.com"}
	subject := "Test mail"
	message := "Hello"

	err := sendMail(conf, to, cc, subject, message)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Mail sent!")
}
