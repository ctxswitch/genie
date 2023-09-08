package http

import (
	"bytes"
	"io"
	"net/http"
	"sync"
	"time"

	"ctx.sh/genie/pkg/resources"
	"ctx.sh/genie/pkg/variables"
	"ctx.sh/strata"
	"github.com/go-logr/logr"
)

type HTTPOptions struct {
	Logger  logr.Logger
	Metrics *strata.Metrics
}

type HTTP struct {
	url     string
	headers Headers
	client  http.Client
	timeout time.Duration
	method  string
	// backoff time.Duration
	logger  logr.Logger
	metrics *strata.Metrics

	resources *resources.Resources
	variables *variables.Variables

	sendChan chan []byte
	stopOnce sync.Once
}

func New(cfg Config, opts *HTTPOptions) *HTTP {
	return &HTTP{
		url:      cfg.Url,
		method:   cfg.Method,
		headers:  newHeaders(cfg.Headers),
		sendChan: make(chan []byte),
		logger:   opts.Logger,
		metrics:  opts.Metrics,
	}
}

func (h *HTTP) Init() error {
	h.sendChan = make(chan []byte)

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

func (h *HTTP) Start() {
	h.start()
}

func (h *HTTP) start() {
	for data := range h.sendChan {
		h.send(data)
	}
}

// TODO: this is still going to be blocking. I need to make this async.
func (h *HTTP) send(data []byte) error {
	req, err := http.NewRequest(h.method, h.url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	for k, v := range h.headers {
		req.Header.Set(k, v.Execute(h.resources, h.variables))
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
		close(h.sendChan)
	})
}

func (h *HTTP) SendChannel() chan<- []byte {
	return h.sendChan
}
