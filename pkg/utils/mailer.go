package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"mime"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
)

type EncryptionType string

const (
	EncryptionNone     EncryptionType = "none"
	EncryptionStartTLS EncryptionType = "starttls"
	EncryptionSSL      EncryptionType = "ssl"
	EncryptionTLS      EncryptionType = "tls"
)

type EmailParams struct {
	FromName   string
	FromEmail  string
	Password   string
	Host       string
	Port       string
	Encryption EncryptionType // none | starttls | ssl | tls
	ToName     string
	ToEmail    string
	Subject    string
	Template   string
	Data       any
	Attachment string // optional path file
}

func SendEmail(p EmailParams) error {
	tmpl, err := template.ParseFiles(p.Template)
	if err != nil {
		return err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, p.Data); err != nil {
		return err
	}

	message, err := buildEmailMessage(p, body.String())
	if err != nil {
		return err
	}

	encryption := strings.ToLower(strings.TrimSpace(string(p.Encryption)))
	switch encryption {
	case "", "starttls":
		return sendWithStartTLS(p, message)
	case "ssl", "tls":
		return sendWithTLS(p, message)
	case "none":
		return sendWithoutTLS(p, message)
	default:
		return fmt.Errorf("unsupported encryption: %s", p.Encryption)
	}
}

func buildEmailMessage(p EmailParams, htmlContent string) ([]byte, error) {
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
		encodeRFC2047(p.Subject),
		boundary,
	)

	if _, err := message.WriteString(headers); err != nil {
		return nil, err
	}

	htmlHeader := textproto.MIMEHeader{}
	htmlHeader.Set("Content-Type", `text/html; charset="UTF-8"`)

	part, err := writer.CreatePart(htmlHeader)
	if err != nil {
		return nil, err
	}

	if _, err := part.Write([]byte(htmlContent)); err != nil {
		return nil, err
	}

	if p.Attachment != "" {
		file, err := os.Open(p.Attachment)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		filename := filepath.Base(p.Attachment)
		encodedFilename := mime.QEncoding.Encode("UTF-8", filename)

		attachmentHeader := textproto.MIMEHeader{}
		attachmentHeader.Set("Content-Type", "application/octet-stream")
		attachmentHeader.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, encodedFilename))
		attachmentHeader.Set("Content-Transfer-Encoding", "base64")

		part, err := writer.CreatePart(attachmentHeader)
		if err != nil {
			return nil, err
		}

		buf := new(bytes.Buffer)
		encoder := base64.NewEncoder(base64.StdEncoding, newBase64LineWriter(buf))

		if _, err := io.Copy(encoder, file); err != nil {
			encoder.Close()
			return nil, err
		}
		if err := encoder.Close(); err != nil {
			return nil, err
		}

		if _, err := part.Write(buf.Bytes()); err != nil {
			return nil, err
		}
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	return message.Bytes(), nil
}

func sendWithStartTLS(p EmailParams, message []byte) error {
	addr := fmt.Sprintf("%s:%s", p.Host, p.Port)

	client, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer client.Close()

	if err := client.StartTLS(&tls.Config{
		ServerName: p.Host,
	}); err != nil {
		return err
	}

	auth := smtp.PlainAuth("", p.FromEmail, p.Password, p.Host)
	if err := client.Auth(auth); err != nil {
		return err
	}

	if err := client.Mail(p.FromEmail); err != nil {
		return err
	}
	if err := client.Rcpt(p.ToEmail); err != nil {
		return err
	}

	wc, err := client.Data()
	if err != nil {
		return err
	}

	if _, err := wc.Write(message); err != nil {
		wc.Close()
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}

	return client.Quit()
}

func sendWithTLS(p EmailParams, message []byte) error {
	addr := fmt.Sprintf("%s:%s", p.Host, p.Port)

	conn, err := tls.Dial("tcp", addr, &tls.Config{
		ServerName: p.Host,
	})
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, p.Host)
	if err != nil {
		return err
	}
	defer client.Close()

	auth := smtp.PlainAuth("", p.FromEmail, p.Password, p.Host)
	if err := client.Auth(auth); err != nil {
		return err
	}

	if err := client.Mail(p.FromEmail); err != nil {
		return err
	}
	if err := client.Rcpt(p.ToEmail); err != nil {
		return err
	}

	wc, err := client.Data()
	if err != nil {
		return err
	}

	if _, err := wc.Write(message); err != nil {
		wc.Close()
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}

	return client.Quit()
}

func sendWithoutTLS(p EmailParams, message []byte) error {
	addr := fmt.Sprintf("%s:%s", p.Host, p.Port)

	client, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer client.Close()

	auth := smtp.PlainAuth("", p.FromEmail, p.Password, p.Host)
	if err := client.Auth(auth); err != nil {
		return err
	}

	if err := client.Mail(p.FromEmail); err != nil {
		return err
	}
	if err := client.Rcpt(p.ToEmail); err != nil {
		return err
	}

	wc, err := client.Data()
	if err != nil {
		return err
	}

	if _, err := wc.Write(message); err != nil {
		wc.Close()
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}

	return client.Quit()
}

func encodeRFC2047(s string) string {
	return mime.QEncoding.Encode("UTF-8", s)
}

type base64LineWriter struct {
	w     io.Writer
	count int
}

func newBase64LineWriter(w io.Writer) io.Writer {
	return &base64LineWriter{w: w}
}

func (lw *base64LineWriter) Write(p []byte) (int, error) {
	written := 0
	for _, b := range p {
		if lw.count == 76 {
			if _, err := lw.w.Write([]byte("\r\n")); err != nil {
				return written, err
			}
			lw.count = 0
		}
		if _, err := lw.w.Write([]byte{b}); err != nil {
			return written, err
		}
		lw.count++
		written++
	}
	return written, nil
}
