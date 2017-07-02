package email

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/smtp"
)

// accountCredentials provides data to establish SMTP-connection.
type accountCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
}

var config accountCredentials

func init() {
	parseConfig()
}

func parseConfig() {
	var err error

	file, err := ioutil.ReadFile("email_config.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatal(err)
	}
}

// SendAlert sends alert e-mail via SMTP to given recipient with proveded text.
func SendAlert(to, text string) error {
	var err error

	//Set up authentication information.
	auth := smtp.PlainAuth("", config.Email, config.Password, config.Host)

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	recipients := []string{to}
	subject := "Subject: System Alert\n"
	msg := []byte(subject + mime + text)

	addr := config.Host + ":" + config.Port

	err = smtp.SendMail(addr, auth, config.Email, recipients, msg)
	if err != nil {
		return err
	}

	return nil
}
