package service

import (
    "bytes"
    "io/ioutil"
    "context"
    "net/http"
    "github.com/op/go-logging"
)

type IHttpClient interface {
    //PostMeasurement(ctx context.Context, thing *model.Thing, value string)
    PostString(ctx context.Context, url, body string, username *string, password *string)
}

type HttpClient struct {
    client *http.Client
}

func NewHttpClient() IHttpClient {

    httpClient := &HttpClient{}
    httpClient.client = &http.Client{}

    return httpClient
}

func (c *HttpClient) PostString(ctx context.Context, url, body string, username *string, password *string) {
    ctx.Value("log").(*logging.Logger).Debugf("Http POST to %s", url)

    req, err := http.NewRequest("POST", url, bytes.NewReader([]byte(body)))
    if username != nil && password != nil {
        req.SetBasicAuth(*username, *password)
    }
    res, err := c.client.Do(req)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Http Post failed (%s)", err.Error())
        return
    }
    response, err := ioutil.ReadAll(res.Body)
    if err != nil {
        ctx.Value("log").(*logging.Logger).Errorf("Http Post failed (%s)", err.Error())
    }

    res.Body.Close()

    ctx.Value("log").(*logging.Logger).Debugf("Http post response status code: %d", res.StatusCode)
    ctx.Value("log").(*logging.Logger).Debugf("Http post response body: %s", response)
}
