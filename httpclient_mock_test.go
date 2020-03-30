package main_test

import (
    "github.com/op/go-logging"
)

type httpClientMockCall struct {
    Url string
    Body string
    Username *string
    Password *string
}

// implements IMqtt interface
type HttpClientMock struct {
    Log *logging.Logger
    Calls []httpClientMockCall
}

func (c *HttpClientMock) PostString(url, body string, username *string, password *string) {

    c.Log.Debugf("Mock Http Client - POST to %s", url)
    c.Calls = append(c.Calls, httpClientMockCall{url, body, username, password})
}
