package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"strconv"
	"time"
)

func SendEmail(to string, subject string, data map[string]any, templateFiles ...string) error {
	host := os.Getenv("MAIL_HOST")
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	username := os.Getenv("MAIL_USERNAME")
	password := os.Getenv("MAIL_PASSWORD")
	fromAddress := os.Getenv("MAIL_FROM_ADDRESS")
	fromName := os.Getenv("MAIL_FROM_NAME")

	if data == nil {
		data = make(map[string]any)
	}

	// Auto-inject default mail data
	data["Title"] = subject
	data["AuthorName"] = "Nova Ardiansyah"
	data["AuthorFirstName"] = "Nova"
	data["Year"] = time.Now().Year()

	var body bytes.Buffer
	t, err := template.ParseFiles(templateFiles...)
	if err != nil {
		return err
	}

	if err := t.Execute(&body, data); err != nil {
		return err
	}

	auth := smtp.PlainAuth("", username, password, host)
	addr := fmt.Sprintf("%s:%d", host, port)

	header := make(map[string]string)
	header["From"] = fmt.Sprintf("%s <%s>", fromName, fromAddress)
	header["To"] = to
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"UTF-8\""

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body.String()

	return smtp.SendMail(addr, auth, fromAddress, []string{to}, []byte(message))
}
