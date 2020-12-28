package main_test

import (
    "github.com/op/go-logging"
)

type mailClientMockCall struct {
    Subject string
    From string
    To []string
    Message string
}

// implements IMqtt interface
type MailClientMock struct {
    Log *logging.Logger
    Calls []mailClientMockCall
}

func (c *MailClientMock) SendMail(subject, from string, to []string, message string) error {
    c.Log.Debugf("Mock Mail Client Called")
    c.Log.Debugf(" - mail subject %s", subject)
    c.Log.Debugf(" - mail from %s", from)
    c.Log.Debugf(" - mail from %v", to)
    c.Log.Debugf(" - mail body %s", message)
    c.Calls = append(c.Calls, mailClientMockCall{subject, from, to, message})
    return nil
}
