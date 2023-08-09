package main

import (
	"email-sender/helper"
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"os"
	// "path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Email string
	Name string
	Password string
	Host string
	Port int
}

type Attachment struct {
	Filename string
	Data     []byte
	Inline   bool
}

type Message struct {
    To[] string
    Subject string
    Message string
    Attachments map[string]*Attachment
}

func getEnv(key string) string {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading.env file")
    }
    return os.Getenv(key)
}

func (c *Config) Load() *Config {
    c.Email = getEnv("CONFIG_AUTH_EMAIL")
	c.Name = getEnv("CONFIG_SENDER_NAME")
	c.Password = getEnv("CONFIG_AUTH_PASSWORD")
	c.Host = getEnv("CONFIG_SMTP_HOST")
	c.Port, _ = strconv.Atoi(getEnv("CONFIG_SMTP_PORT"))
    return c
}

func (m *Message) New(to []string, subject string, message string, attachments map[string]*Attachment) *Message {
    m.To = to
    m.Subject = subject
    m.Message = message
    m.Attachments = attachments
    return m
}

func Send(conf Config, email Message) error {
	body := "From: " + conf.Name + "\n" +
		    "To: " + strings.Join(email.To, ",") + "\n" +
		    "Subject: " + email.Subject + "\n\n" +
		    email.Message

	auth := smtp.PlainAuth("", conf.Email, conf.Password, conf.Host)
	smtpAddr := fmt.Sprintf("%s:%d", conf.Host, conf.Port)

	err := smtp.SendMail(smtpAddr, auth, conf.Email, email.To, []byte(body))
	if err != nil {
		return err
	}

	return nil
}


func loadFromJsonFile(relativeFilePath string) helper.BodyJson {
    absoluteFilePath := helper.GenerateAbsolutePath(relativeFilePath)
    
    fileContentJson, err := helper.ReadFile(absoluteFilePath)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Print(string(fileContentJson))

    var body helper.BodyJson
    err = json.Unmarshal(fileContentJson, &body)
    if err != nil {
        log.Fatal(err)
    }

    return body
}

func main() {
    relativeFilePath := "resource/body.json"
    j := loadFromJsonFile(relativeFilePath)

    message := strings.ReplaceAll(strings.ReplaceAll(j.Message, "[nama]", "Test"), "[perusahaan]", "PT Djarum")
    to := []string{"ptadityamahendrap@gmail.com", "whatupbiatch69@gmail.com"}
    subject := j.Subject
    
    var m Message
    var conf Config 
    
    conf.Load()
    m.New(to, subject, message, nil)

	err := Send(conf, m)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Mail sent!")
}
