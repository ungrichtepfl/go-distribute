package distribute

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"net/smtp"
	"path/filepath"
	"strings"
)

func SendEmails(config Config, name_email_to_documents map[NameEmail][]DocumentPath) map[NameEmail]error {
	name_email_to_error := make(map[NameEmail]error)
	for name_email, documents := range name_email_to_documents {
		err := SendEmail(config, name_email.Email, documents)
		if err != nil {
			name_email_to_error[name_email] = err
		}
	}
	return name_email_to_error

}

func SendEmail(config Config, email string, documents []DocumentPath) error {

	addr := fmt.Sprintf("%s:%s", config.EmailConfig.SMTPHost, config.EmailConfig.SMTPPort)
	message := config.AsEmailMessage(email, documents)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, config.EmailConfig.SMTPHost)
	if err != nil {
		return err
	}

	tlsconfig := &tls.Config{
		ServerName: config.EmailConfig.SMTPHost,
	}

	if err = c.StartTLS(tlsconfig); err != nil {
		return err
	}
	auth := LoginAuth(config.EmailConfig.Username, config.EmailConfig.Password)
	if err = c.Auth(auth); err != nil {
		return err
	}

	// TODO: Check if the above works also for gmail etc. or only the below is needed
	// auth := smtp.PlainAuth("", config.EmailConfig.Username, config.EmailConfig.Password, config.EmailConfig.SMTPHost)

	log.Println("Start sending email to", email)
	err = smtp.SendMail(addr, auth, config.EmailConfig.SenderEmail, message.To, message.ToBytes())
	log.Println("Sent email to", email)
	return err
}

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unknown from server")
		}
	}
	return nil, nil
}

func (m *Message) AttachFile(src string) error {
	b, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	_, fileName := filepath.Split(src)
	m.Attachments[fileName] = b
	return nil
}

func (m *Message) ToBytes() []byte {
	buf := bytes.NewBuffer(nil)
	withAttachments := len(m.Attachments) > 0
	buf.WriteString(fmt.Sprintf("Subject: %s\n", m.Subject))
	buf.WriteString(fmt.Sprintf("To: %s\n", strings.Join(m.To, ",")))
	if len(m.CC) > 0 {
		buf.WriteString(fmt.Sprintf("Cc: %s\n", strings.Join(m.CC, ",")))
	}

	if len(m.BCC) > 0 {
		buf.WriteString(fmt.Sprintf("Bcc: %s\n", strings.Join(m.BCC, ",")))
	}

	buf.WriteString("MIME-Version: 1.0\n")
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()
	if withAttachments {
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n", boundary))
		buf.WriteString(fmt.Sprintf("--%s\n", boundary))
	} else {
		buf.WriteString("Content-Type: text/plain; charset=utf-8\n")
	}

	buf.WriteString(m.Body)
	if withAttachments {
		for k, v := range m.Attachments {
			buf.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s\n", http.DetectContentType(v)))
			buf.WriteString("Content-Transfer-Encoding: base64\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\n", k))

			b := make([]byte, base64.StdEncoding.EncodedLen(len(v)))
			base64.StdEncoding.Encode(b, v)
			buf.Write(b)
			buf.WriteString(fmt.Sprintf("\n--%s", boundary))
		}

		buf.WriteString("--")
	}

	return buf.Bytes()
}
