package tinify

import "time"

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
