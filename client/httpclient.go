package main

import (
	"fmt"
	"net/http"
	"bytes"
	"io/ioutil"
)

type RcpRestClient struct {
	Conf Conf
}

func NewRcpRestClient(c Conf) *RcpRestClient {
  return &RcpRestClient{Conf: c}
}

func (c *RcpRestClient) DoPut(txt string) error {
	req, err := http.NewRequest("PUT", c.Conf.Host, bytes.NewBufferString(txt))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	return err
}

func (c *RcpRestClient) DoGet() (string, error) {
	req, err := http.NewRequest("GET", c.Conf.Host, nil) 
	if err != nil {
		return "", err
	}

	return c.doRequest(req)
}

func (c *RcpRestClient) DoDelete() error {
	req, err := http.NewRequest("DELETE", c.Conf.Host, nil) 
	if err != nil {
		return err
	}
	_, err = c.doRequest(req)
	return err
}

func (c *RcpRestClient) doRequest(req *http.Request) (string, error) {
	req.Header.Add("Authorization", "Bearer "+c.Conf.Secret)
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if 200 != resp.StatusCode {
		return "", fmt.Errorf("%s", body)
	}

	return string(body), nil
}
