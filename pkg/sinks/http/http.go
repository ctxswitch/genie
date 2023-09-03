package http

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"sync"
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
	sendChan chan []byte
	stopChan chan struct{}
	stopOnce sync.Once
}

func New() *HTTP {
	return &HTTP{
		url:      DefaultHttpUrl,
		method:   DefaultMethod,
		headers:  make(map[string]any, 0),
		sendChan: make(chan []byte),
		stopChan: make(chan struct{}),
	}
}

func (h *HTTP) Init() error {
	h.sendChan = make(chan []byte)
	h.stopChan = make(chan struct{})

	h.client = http.Client{
		Timeout: h.timeout,
		Transport: &http.Transport{
			// TODO: make configurable for connecction pooling
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableKeepAlives:  false,
			DisableCompression: false,
		},
	}

	return nil
}

func (h *HTTP) Start(ctx context.Context) {
	h.start(ctx)
}

func (h *HTTP) start(ctx context.Context) {
	for {
		select {
		case <-h.stopChan:
			return
		case <-ctx.Done():
			return
		case d := <-h.sendChan:
			h.send(d)
		}
	}
}

// TODO: this is still going to be blocking. I need to make this async.
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
			// log and then remove the invalid header
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

func (h *HTTP) Stop() {
	h.stopOnce.Do(func() {
		close(h.stopChan)
	})
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

func (h *HTTP) SendChannel() chan<- []byte {
	return h.sendChan
}
