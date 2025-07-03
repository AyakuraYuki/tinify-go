package tinify

import (
	"errors"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
)

type ResizeMethod string

func (rm ResizeMethod) String() string { return string(rm) }

const (
	ResizeMethodScale ResizeMethod = "scale"
	ResizeMethodFit   ResizeMethod = "fit"
	ResizeMethodCover ResizeMethod = "cover"
)

type ResizeOption struct {
	Method ResizeMethod `json:"method"`
	Width  int64        `json:"width"`
	Height int64        `json:"height"`
}

type Source struct {
	url      string
	commands map[string]any
}

func newSource(url string, commands map[string]any) *Source {
	s := &Source{
		url:      url,
		commands: make(map[string]any),
	}
	if len(commands) > 0 {
		s.commands = commands
	}
	return s
}

func fromResponse(rsp *resty.Response) (source *Source, err error) {
	location := rsp.Header().Get("Location")
	source = newSource(location, nil)
	return source, nil
}

func (c *Client) FromFile(path string) (source *Source, err error) {
	buffer, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return c.FromBuffer(buffer)
}

func (c *Client) FromBuffer(buffer []byte) (source *Source, err error) {
	rsp, err := c.request(methodPOST, "/shrink", buffer)
	if err != nil {
		return nil, err
	}
	source, err = fromResponse(rsp)
	return
}

func (c *Client) FromURL(u string) (source *Source, err error) {
	url := strings.TrimSpace(u)
	if url == "" {
		return nil, errors.New("url is required")
	}

	body := map[string]any{
		"source": map[string]any{
			"url": url,
		},
	}

	rsp, err := c.request(methodPOST, "/shrink", body)
	if err != nil {
		return nil, err
	}
	source, err = fromResponse(rsp)
	return
}

func (c *Client) ToFile(source *Source, dst string) (err error) {
	result, err := c.toResult(source)
	if err != nil {
		return err
	}
	return result.ToFile(dst)
}

func (c *Client) Resize(source *Source, option *ResizeOption) (err error) {
	if source == nil {
		return errors.New("source is required")
	}
	if option == nil {
		return errors.New("option is required")
	}

	source.commands["resize"] = option
	return nil
}

func (c *Client) toResult(source *Source) (result *Result, err error) {
	if source == nil || len(source.url) == 0 {
		return nil, errors.New("no valid source")
	}

	rsp, err := c.request(methodGET, source.url, source.commands)
	if err != nil {
		return nil, err
	}

	result = NewResult(rsp.Header(), rsp.Body())
	return result, nil
}
