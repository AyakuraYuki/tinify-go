package tinify

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	baseURL = "https://api.tinify.com"
)

type ClientOption func(client *Client)

// WithProxy configures the client to use the specified proxy for making requests,
//
//	WithProxy("http://proxyserver:8888")
//
// or you could also set proxy via environment variable, refer to [http.ProxyFromEnvironment]
func WithProxy(proxy string) ClientOption {
	return func(client *Client) {
		client.proxy = proxy
	}
}

// WithAppIdentifier sets the app identifier for the client, will be appended to the User-Agent.
func WithAppIdentifier(appIdentifier string) ClientOption {
	return func(client *Client) {
		client.appIdentifier = appIdentifier
	}
}

// WithRetryCount sets the times of retries for each request, no retry if set to 0.
func WithRetryCount(retry int) ClientOption {
	return func(client *Client) {
		client.retry = retry
	}
}

// WithRetryWaitTime sets the wait time for sleep before retrying request.
func WithRetryWaitTime(duration time.Duration) ClientOption {
	return func(client *Client) {
		client.retryWaitTime = duration
	}
}

type Client struct {
	key           string
	appIdentifier string
	retry         int
	retryWaitTime time.Duration
	proxy         string

	client           *resty.Client
	compressionCount int
}

// NewClient creates a new Client instance with the provided API key and optional configuration options.
func NewClient(key string, opts ...ClientOption) *Client {
	c := &Client{
		key:           key,
		appIdentifier: "",
		retry:         1,
		retryWaitTime: 500 * time.Millisecond,
		proxy:         "",
	}

	for _, opt := range opts {
		opt(c)
	}

	userAgent := fmt.Sprintf(
		"Tinify/%s Golang/%s (%s %s)",
		Version,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH)
	if c.appIdentifier != "" {
		userAgent += " " + c.appIdentifier
	}

	cli := resty.New()
	cli.SetRetryCount(c.retry)
	cli.SetRetryWaitTime(c.retryWaitTime)
	cli.SetBasicAuth("api", c.key)
	cli.SetHeader("User-Agent", userAgent)

	if c.proxy != "" {
		cli.SetProxy(c.proxy)
	}

	c.client = cli

	return c
}

type Method string

const (
	methodGET  Method = http.MethodGet
	methodPOST Method = http.MethodPost
)

func (c *Client) request(method Method, endpoint string, body any) (rsp *resty.Response, err error) {
	req := c.client.R().SetBody(body)

	url := endpoint
	if !strings.HasPrefix(url, "http") {
		url = baseURL + endpoint
	}

	if body != nil {
		switch body.(type) {
		case map[string]any:
			req.SetHeader("Content-Type", "application/json")
		}
	}

	switch method {
	case methodGET:
		rsp, err = req.Get(url)
	case methodPOST:
		rsp, err = req.Post(url)
	default:
		err = errors.New("unsupported method")
	}
	if err != nil {
		return nil, err
	}

	if status := rsp.StatusCode(); status >= http.StatusOK && status < http.StatusMultipleChoices {
		return rsp, nil
	}

	ex := &ErrorData{}
	if err = json.Unmarshal(rsp.Body(), &ex); err != nil {
		return nil, err
	}

	return nil, ex
}
