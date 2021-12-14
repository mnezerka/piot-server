package main

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/op/go-logging"
)

type IHttpClient interface {
	//PostMeasurement(ctx context.Context, thing *Thing, value string)
	PostString(url, body string, username *string, password *string)
}

type HttpClient struct {
	log    *logging.Logger
	client *http.Client
}

func NewHttpClient(log *logging.Logger) IHttpClient {

	httpClient := &HttpClient{log: log}
	httpClient.client = &http.Client{}

	return httpClient
}

func (c *HttpClient) PostString(url, body string, username *string, password *string) {
	c.log.Debugf("Http POST to %s", url)

	req, err := http.NewRequest("POST", url, bytes.NewReader([]byte(body)))
	if err != nil {
		c.log.Errorf("Http Post failed (%s)", err.Error())
		return
	}

	if username != nil && password != nil {
		req.SetBasicAuth(*username, *password)
	}
	res, err := c.client.Do(req)
	if err != nil {
		c.log.Errorf("Http Post failed (%s)", err.Error())
		return
	}
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.log.Errorf("Http Post failed (%s)", err.Error())
	}

	res.Body.Close()

	c.log.Debugf("Http post response status code: %d", res.StatusCode)
	c.log.Debugf("Http post response body: %s", response)
}
