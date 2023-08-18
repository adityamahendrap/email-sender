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

	"github.com/360EntSecGroup-Skylar/excelize"
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


func ReadFromJsonFile(relativeFilePath string) helper.BodyJson {
    absoluteFilePath := helper.GenerateAbsolutePath(relativeFilePath)
    
    fileContentJson, err := helper.ReadFile(absoluteFilePath)
    if err != nil {
        log.Fatal(err)
    }
    // fmt.Print(string(fileContentJson))

    var body helper.BodyJson
    err = json.Unmarshal(fileContentJson, &body)
    if err != nil {
        log.Fatal(err)
    }

    return body
}

type M map[string]interface{}

func ReadFromExcelFile(relativePath string, sheetName string, col string, fromRow int, toRow int) []string {
    xlsx, err := excelize.OpenFile(relativePath)
    if err != nil {
        log.Fatal(err.Error())
    }

    var rows []string
    for i := fromRow; i <= toRow; i++ {
        row :=  xlsx.GetCellValue(sheetName, fmt.Sprintf("%s%d", col, i))
        rows = append(rows, row)
    }

    // fmt.Printf("%v \n", rows)

    return rows
}

func main() {
    // names := ReadFromExcelFile("./resource/data.xlsx","Lomba - Ide Bisnis", "C", 3, 45)
    // emails := ReadFromExcelFile("./resource/data.xlsx", "Lomba - Ide Bisnis", "D", 3, 45)
    j := ReadFromJsonFile("resource/body.json")

    // for i := 0; i < len(names); i++ {
    //     fmt.Printf("%s %s\n", names[i], emails[i])
    // }

    names := []string{"Test", "Test2"}
    emails := []string{"whatupbiatch69@gmail.com", "ptadityamahendrap@gmail.com"}

    done := make(chan bool) 
    
    var sendedCount int = 0
    for i := 0; i < len(emails); i++ {
        go func(name, email string) {
            defer func() {
                done <- true 
            }()
            
            message := strings.ReplaceAll(strings.ReplaceAll(j.Message, "[nama]", "Test"), "[perusahaan]", name)
            to := []string{email}
            subject := j.Subject

            var m Message
            var conf Config

            conf.Load()
            m.New(to, subject, message, nil)

            log.Printf("Sending mail to %s...\n", email)
            err := Send(conf, m)
            if err != nil {
                log.Printf("Error sending mail to %s: %s\n", email, err.Error())
            } else {
                log.Printf("Mail sent to %s\n", email)
            }
        }(names[i], emails[i])
    }

    for i := 0; i < len(emails); i++ {
        <-done
        sendedCount++
    }

    log.Println("Done! Total sended email:", sendedCount)
}
