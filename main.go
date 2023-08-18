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
    "bytes"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/joho/godotenv"
    "encoding/base64"
)

type Config struct {
	Email string
	Name string
	Password string
	Host string
	Port int
}

type Attachment struct {
    FileName string
    Data     []byte
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

func AttachFiles(message *Message, filenames []string) error {
    attachments := make(map[string]*Attachment)

    for _, filename := range filenames {
        data, err := os.ReadFile(filename)
        if err != nil {
            return err
        }

        attachments[filename] = &Attachment{
            FileName: filename,
            Data:     data,
        }
    }

    message.Attachments = attachments
    return nil
}

func Send(conf Config, email Message) error {
    boundary := "boundary123"
    
    var buf bytes.Buffer

    // Build the email headers
    buf.WriteString("From: " + conf.Name + "\n")
    buf.WriteString("To: " + strings.Join(email.To, ",") + "\n")
    buf.WriteString("Subject: " + email.Subject + "\n")
    buf.WriteString("MIME-Version: 1.0\n")
    buf.WriteString("Content-Type: multipart/mixed; boundary=" + boundary + "\n\n")

    // Start the message body
    buf.WriteString("--" + boundary + "\n")
    buf.WriteString("Content-Type: text/plain; charset=utf-8\n\n")
    buf.WriteString(email.Message + "\n\n")

    // Attachments for jpg file
    for _, attachment := range email.Attachments {
        buf.WriteString("--" + boundary + "\n")
        buf.WriteString("Content-Type: image/jpeg\n")
        buf.WriteString("Content-Disposition: attachment; filename=\"" + attachment.FileName + "\"\n")
        buf.WriteString("Content-Transfer-Encoding: base64\n\n")
        
        base64Data := base64.StdEncoding.EncodeToString(attachment.Data)
        buf.WriteString(base64Data + "\n")
    }

    buf.WriteString("--" + boundary + "--\n")

    auth := smtp.PlainAuth("", conf.Email, conf.Password, conf.Host)
    smtpAddr := fmt.Sprintf("%s:%d", conf.Host, conf.Port)

    err := smtp.SendMail(smtpAddr, auth, conf.Email, email.To, buf.Bytes())
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
    names := ReadFromExcelFile("./resource/data.xlsx","Lomba - Ide Bisnis", "C", 3, 45)
    emails := ReadFromExcelFile("./resource/data.xlsx", "Lomba - Ide Bisnis", "D", 3, 45)
    // emails := []string{"ptadityamahendrap@gmail.com"}
    j := ReadFromJsonFile("resource/itcc.json")

    for i := 0; i < len(emails); i++ {
        fmt.Printf("%s - %s\n", names[i], emails[i])
    }

    return
    
    message := j.Message
    subject := j.Subject

    done := make(chan bool) 
    
    var sendedCount int = 0

    attachments := []string{"test.png", "PAMFLET ITCC.jpg"}
    var attachmentMap = make(map[string]*Attachment)
    for _, attachment := range attachments {
        data, err := os.ReadFile(attachment)
        if err != nil {
            log.Printf("Error reading attachment %s: %s\n", attachment, err.Error())
            continue
        }
        attachmentMap[attachment] = &Attachment{
            FileName: attachment,
            Data:     data,
        }
    }

    for i := 0; i < len(emails); i++ {
        go func(name, email string) {
            defer func() {
                done <- true 
            }()
            
            // message := strings.ReplaceAll(strings.ReplaceAll(j.Message, "[nama]", "Test"), "[perusahaan]", name)
            to := []string{email}
            // subject := j.Subject

            var m Message
            var conf Config

            conf.Load()
            m.New(to, subject, message, nil)

            err := AttachFiles(&m, attachments) 
            if err != nil {
                log.Printf("Error attaching files to email: %s\n", err.Error())
                return
            }

            log.Printf("Sending mail to %s...\n", email)
            err = Send(conf, m)
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
