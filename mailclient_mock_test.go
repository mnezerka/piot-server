package main_test

import (
    "github.com/op/go-logging"
)

type mailClientMockCall struct {
    From string
    To []string
    Message string
}

// implements IMqtt interface
type MailClientMock struct {
    Log *logging.Logger
    Calls []mailClientMockCall
}

func (c *MailClientMock) SendMail(from string, to []string, message string) error {
    c.Log.Debugf("Mock Mail Client - mail from %s", from)
    c.Calls = append(c.Calls, mailClientMockCall{from, to, message})
    return nil
}
