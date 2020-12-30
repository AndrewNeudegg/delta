package e2e

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type SinkClient struct {
	Addr string
}

func (s SinkClient) Send(URI string, headers map[string][]string, content []byte) (string, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", s.Addr, URI), bytes.NewBuffer(content))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
