package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"net/smtp"
	"os"
	"path/filepath"
)

type EmailParams struct {
	FromName   string
	FromEmail  string
	Password   string
	ToName     string
	ToEmail    string
	Subject    string
	Template   string
	Data       any
	Attachment string // optional path file
}

func SendEmail(p EmailParams) error {

	// render template
	tmpl, err := template.ParseFiles(p.Template)
	if err != nil {
		return err
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, p.Data)
	if err != nil {
		return err
	}

	htmlContent := body.String()

	var message bytes.Buffer
	writer := multipart.NewWriter(&message)

	boundary := writer.Boundary()

	headers := fmt.Sprintf(
		"From: %s <%s>\r\n"+
			"To: %s <%s>\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: multipart/mixed; boundary=%s\r\n\r\n",
		p.FromName,
		p.FromEmail,
		p.ToName,
		p.ToEmail,
		p.Subject,
		boundary,
	)

	message.WriteString(headers)

	// HTML body
	part, _ := writer.CreatePart(map[string][]string{
		"Content-Type": {"text/html; charset=UTF-8"},
	})

	part.Write([]byte(htmlContent))

	// attachment optional
	if p.Attachment != "" {

		file, err := os.Open(p.Attachment)
		if err != nil {
			return err
		}
		defer file.Close()

		filename := filepath.Base(p.Attachment)

		part, _ := writer.CreatePart(map[string][]string{
			"Content-Type":              {"application/octet-stream"},
			"Content-Disposition":       {fmt.Sprintf(`attachment; filename="%s"`, filename)},
			"Content-Transfer-Encoding": {"base64"},
		})

		buf := new(bytes.Buffer)
		encoder := base64.NewEncoder(base64.StdEncoding, buf)

		io.Copy(encoder, file)
		encoder.Close()

		part.Write(buf.Bytes())
	}

	writer.Close()

	auth := smtp.PlainAuth("", p.FromEmail, p.Password, "smtp.gmail.com")

	return smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		p.FromEmail,
		[]string{p.ToEmail},
		message.Bytes(),
	)
}
