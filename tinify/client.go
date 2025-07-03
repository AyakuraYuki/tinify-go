package tinify

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	baseURL = "https://api.tinify.com"
)

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

func (c *Client) CompressionCount() int {
	return c.compressionCount
}

func (c *Client) request(method string, endpoint string, body any) (rsp *resty.Response, err error) {
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
	case http.MethodGet:
		rsp, err = req.Get(url)
	case http.MethodPost:
		rsp, err = req.Post(url)
	default:
		err = errors.New("unsupported method")
	}
	if err != nil {
		return nil, err
	}

	if compressionCount := rsp.Header().Get("Compression-Count"); compressionCount != "" {
		cc, _ := strconv.Atoi(compressionCount)
		if cc >= c.compressionCount {
			c.compressionCount = cc
		}
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
