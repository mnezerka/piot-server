package main

import (
	"fmt"
	"net/smtp"
	"piot-server/config"
	"strings"

	"github.com/op/go-logging"
)

type IMailClient interface {
	SendMail(subject, from string, to []string, message string) error
}

type MailClient struct {
	log    *logging.Logger
	params *config.Parameters
	//client *http.Client
}

func NewMailClient(log *logging.Logger, params *config.Parameters) IMailClient {
	mailClient := &MailClient{log: log, params: params}
	//httpClient.client = &http.Client{}
	return mailClient
}

/*
Commented out since it was not possbile to disable authentication, which is
not needed if sending to docker swarm/compose smtp relay host

func (c *MailClient) SendMail(from string, to []string, message string) error {
    c.log.Debugf("SendMail from %s to %v", from, to)
    c.log.Debugf("- smtp host: %s", c.params.SmtpHost)
    c.log.Debugf("- smtp port: %d", c.params.SmtpPort)
    c.log.Debugf("- smtp user: %s", c.params.SmtpUser)

    if len(from) == 0 {
        return fmt.Errorf("Cannot send mail to %s. From address is empty", to)
    }

    if len(to) == 0 {
        return fmt.Errorf("Cannot send mail, no recepients provided")
    }

    // Authentication.
    auth := smtp.PlainAuth("", c.params.SmtpUser, c.params.SmtpPassword, c.params.SmtpHost)

    msgRaw := []byte("To: " + strings.Join(to, ";") + "\r\n" +
        "Subject: Monitoring!\r\n" +
        "\r\n" +
        message +
        "\r\n")

    c.log.Debugf("message: %s", string(msgRaw))

    // Sending email.
    err := smtp.SendMail(fmt.Sprintf("%s:%d", c.params.SmtpHost, c.params.SmtpPort), auth, from, to, msgRaw)
    if err != nil {
        c.log.Errorf("Cannot send email, error: %v", err)
        return err
    }

    return nil
}
*/
func (c *MailClient) SendMail(subject, from string, to []string, message string) error {

	c.log.Debugf("SendMail from %s to %v", from, to)
	c.log.Debugf("- smtp host: %s", c.params.SmtpHost)
	c.log.Debugf("- smtp port: %d", c.params.SmtpPort)
	c.log.Debugf("- smtp user: %s", c.params.SmtpUser)

	if len(from) == 0 {
		return fmt.Errorf("cannot send mail to %s. From address is empty", to)
	}

	if len(to) == 0 {
		return fmt.Errorf("cannot send mail, no recepients provided")
	}

	r := strings.NewReplacer("\r\n", "", "\r", "", "\n", "", "%0a", "", "%0d", "")

	addr := fmt.Sprintf("%s:%d", c.params.SmtpHost, c.params.SmtpPort)

	client, err := smtp.Dial(addr)
	if err != nil {
		return err
	}

	defer client.Close()
	if err = client.Mail(r.Replace(from)); err != nil {
		return err
	}
	for i := range to {
		to[i] = r.Replace(to[i])
		if err = client.Rcpt(to[i]); err != nil {
			return err
		}
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	msgRaw := "To: " + strings.Join(to, ",") + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
		"\r\n" + message

	c.log.Debugf("message: %s", msgRaw)

	_, err = w.Write([]byte(msgRaw))
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return client.Quit()
}
