package http

import (
	"context"
	"net/http"
	"net/url"
	"time"

	gjson "github.com/sidgwick/gutil/json"
	"github.com/spf13/cast"
)

type Curl struct {
	Url string `json:"url"`
}

type Request struct {
	Data    interface{}   `json:"data"`
	Header  http.Header   `json:"header"`
	Timeout time.Duration `json:"timeout"`
}

func NewCurl(urlStr string) (*Curl, error) {
	curl := &Curl{
		Url: urlStr,
	}

	return curl, nil
}

func (c *Curl) buildUrl(data interface{}) string {
	u, _ := url.Parse(c.Url)

	queryString := BuildGetQueryString(data)
	queryString = MergeQueryString(queryString, u.RawQuery)

	u.RawQuery = queryString
	return u.String()
}

func (c *Curl) Get(ctx context.Context, req *Request, ops ...GetOption) (string, error) {
	withTimeout := WithTimeout(req.Timeout)
	client := MustNewClient(withTimeout)

	u := c.buildUrl(req.Data)
	resp, err := client.Get(ctx, u, req, ops...)
	if err != nil {
		return "", err
	}

	respBody := cast.ToString(resp)
	return respBody, nil
}

func (c *Curl) Post(ctx context.Context, req *Request, ops ...GetOption) (string, error) {
	withTimeout := WithTimeout(req.Timeout)
	client := MustNewClient(withTimeout)

	u := c.buildUrl(nil)

	resp, err := client.Post(ctx, u, req, ops...)
	if err != nil {
		return "", err
	}

	respBody := cast.ToString(resp)
	return respBody, nil
}

func (c *Curl) GetJson(ctx context.Context, req *Request, resp interface{}, ops ...GetOption) error {
	respBody, err := c.Get(ctx, req, ops...)

	err = gjson.LoadData(resp, respBody)
	if err != nil {
		return err
	}

	return nil
}

func (c *Curl) PostJson(ctx context.Context, req *Request, resp interface{}, ops ...GetOption) error {
	respBody, err := c.Post(ctx, req, ops...)

	err = gjson.LoadData(resp, respBody)
	if err != nil {
		return err
	}

	return nil
}
