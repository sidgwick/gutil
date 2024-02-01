package http

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"

	gjson "github.com/sidgwick/gutil/json"
	"github.com/sidgwick/gutil/logger"
)

const (
	defaultTimeout = 3 * time.Second
)

type Client struct {
	client *http.Client
	ops    getoptions
}

type getoptions struct {
	timeout time.Duration
}

type GetOption func(o *getoptions)

func WithTimeout(timeout time.Duration) GetOption {
	return func(o *getoptions) {
		o.timeout = timeout
	}
}

func MustNewClient(opts ...GetOption) *Client {
	options := getoptions{}

	for _, op := range opts {
		op(&options)
	}

	c := Client{
		ops: options,
	}

	timeout := options.timeout
	if timeout == 0 {
		timeout = defaultTimeout
	}

	c.client = &http.Client{
		Timeout: timeout,
	}

	return &c
}

func (c *Client) Get(ctx context.Context, url string, _req *Request, ops ...GetOption) ([]byte, error) {
	options := c.ops
	for _, op := range ops {
		op(&options)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header = _req.Header
	req.Close = true

	resp, err := c.client.Do(req)

	var respStr []byte
	if resp != nil {
		respStr, _ = httputil.DumpResponse(resp, true)
	}

	reqStr, _ := httputil.DumpRequest(req, true)
	logger.Tracef(ctx, "http POST request:\n%v\n\nresponse:\n%v", string(reqStr), string(respStr))

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	return ioutil.ReadAll(resp.Body)
}

func (c *Client) Post(ctx context.Context, url string, _req *Request, ops ...GetOption) ([]byte, error) {
	options := c.ops
	for _, op := range ops {
		op(&options)
	}

	dataStr := gjson.Json(_req.Data)
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(dataStr))
	if err != nil {
		return nil, err
	}

	req.Header = _req.Header
	req.Close = true

	resp, err := c.client.Do(req)

	var respStr []byte
	if resp != nil {
		respStr, _ = httputil.DumpResponse(resp, true)
	}

	reqStr, _ := httputil.DumpRequest(req, true)
	logger.Tracef(ctx, "http POST request:\n%v\n\nresponse:\n%v", string(reqStr), string(respStr))

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	return ioutil.ReadAll(resp.Body)
}
