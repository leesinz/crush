package mail

import (
	"crush/config"
	"crush/utils"
	"fmt"
	"io/ioutil"
	"net/smtp"
	"path/filepath"
	"strings"
	"time"
)

type Email struct {
	SMTPServer string
	SMTPPort   string
	Username   string
	Password   string
	From       string
	To         []string
	Subject    string
	Body       string
}

var (
	cfg        = config.LoadConfig()
	maillogDir = filepath.Join(utils.GetParentPath(), "data", "updateinfo", "log")
)

func NewEmail(smtpServer, smtpPort, username, password, from, subject, body string, to ...string) *Email {
	return &Email{
		SMTPServer: smtpServer,
		SMTPPort:   smtpPort,
		Username:   username,
		Password:   password,
		From:       from,
		To:         to,
		Subject:    subject,
		Body:       body,
	}
}

// Send 方法用于发送邮件
func (email *Email) Send() error {
	message := fmt.Sprintf("Content-Type: text/html; charset=UTF-8\r\n")
	message += fmt.Sprintf("To: %s\r\n", strings.Join(email.To, ", "))
	message += fmt.Sprintf("From: %s\r\n", "vul_monitor")
	message += fmt.Sprintf("Subject: %s\r\n\r\n", email.Subject)
	message += email.Body

	auth := smtp.PlainAuth("", email.Username, email.Password, email.SMTPServer)

	// 发送邮件
	err := smtp.SendMail(email.SMTPServer+":"+email.SMTPPort, auth, email.From, email.To, []byte(message))
	if err != nil {
		return fmt.Errorf("send err: %v", err)
	}

	return nil
}

func Sendmail() {
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	logpath := filepath.Join(maillogDir, yesterday)
	logContent, _ := ioutil.ReadFile(logpath)
	mailcfg := cfg.Email
	email := NewEmail(
		mailcfg.SMTP_SERVER,
		mailcfg.SMTP_PORT,
		mailcfg.Username,
		mailcfg.Password,
		mailcfg.From,
		"Vulnerability Update Monitor",
		string(logContent),
		mailcfg.To...,
	)

	if err := email.Send(); err != nil {
		fmt.Println(err)
		return
	}
	utils.PrintColor("success", "Send mail successfully")
}
