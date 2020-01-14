package service

import (
    "context"
    "github.com/op/go-logging"
    //"piot-server/model"
)

type httpClientMockCall struct {
    Url string
    Body string
    Username *string
    Password *string
}

// implements IMqtt interface
type HttpClientMock struct {
    Calls []httpClientMockCall
}

func (c *HttpClientMock) PostString(ctx context.Context, url, body string, username *string, password *string) {

    ctx.Value("log").(*logging.Logger).Debugf("Mock Http Client - POST to %s", url)

    c.Calls = append(c.Calls, httpClientMockCall{url, body, username, password})
}
