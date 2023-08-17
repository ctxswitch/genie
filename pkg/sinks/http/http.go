package http

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"ctx.sh/genie/pkg/resources"
)

type HTTP struct {
	url     string
	headers map[string]any
	client  http.Client
	// backoff time.Duration
	timeout time.Duration
	method  string
	// logger  *zap.Logger
}

func New() *HTTP {
	return &HTTP{
		url:     DefaultHttpUrl,
		method:  DefaultMethod,
		headers: make(map[string]any, 0),
	}
}

func (h *HTTP) Connect() {
	tr := &http.Transport{
		// TODO: make configurable and add parallel sends
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableKeepAlives:  false,
		DisableCompression: false,
	}
	h.client = http.Client{
		Timeout:   h.timeout,
		Transport: tr,
	}
}

func (h *HTTP) Init() {}

func (h *HTTP) Name() string {
	return "http"
}

func (h *HTTP) Send(data []byte) error {
	// TODO: backoff handling on error
	return h.send(data)
}

func (h *HTTP) Validate() error {
	return nil
}

func (h *HTTP) WithMethod(method string) *HTTP {
	h.method = method
	return h
}

func (h *HTTP) WithURL(url string) *HTTP {
	h.url = url
	return h
}

func (h *HTTP) WithHeader(name string, value any) *HTTP {
	h.headers[name] = value
	return h
}

func (h *HTTP) send(data []byte) error {
	req, err := http.NewRequest(h.method, h.url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	for k, v := range h.headers {
		var val string
		switch h := v.(type) {
		case resources.Resource:
			val = h.Get()
		case string:
			val = h
		default:
			// log
			continue
		}

		req.Header.Set(k, val)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// TODO: log the response body if requested
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}
