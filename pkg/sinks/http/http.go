package http

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"ctx.sh/genie/pkg/resources"
	"ctx.sh/genie/pkg/variables"
	"ctx.sh/strata"
	"github.com/go-logr/logr"
)

// Options are the options for an HTTP sink.
type Options struct {
	Logger    logr.Logger
	Metrics   *strata.Metrics
	Resources *resources.Resources
}

// HTTP is an HTTP sink.
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
	// TODO: look at making this a buffered channel when we allow
	// for multiple sink workers.
	sendChan chan []byte
	stopOnce sync.Once
}

// New returns a new HTTP sink.
func New(cfg Config, opts *Options) *HTTP {
	return &HTTP{
		url:       cfg.URL,
		method:    strings.ToUpper(cfg.Method),
		headers:   newHeaders(cfg.Headers),
		sendChan:  make(chan []byte),
		logger:    opts.Logger,
		metrics:   opts.Metrics,
		resources: opts.Resources,
	}
}

// Init initializes the HTTP sink, setting up the client and send channel
// that events will be sent to.
func (h *HTTP) Init() error {
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

// Start starts the HTTP sink.
func (h *HTTP) Start() {
	h.start()
}

// start is the main loop for the sink that listens for events on the send
// channel.
func (h *HTTP) start() {
	for data := range h.sendChan {
		if err := h.send(data); err != nil {
			h.logger.Error(err, "failed to send data")
		}
	}
}

// send sends the data to the configured URL.
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

// Stop stops the HTTP sink.
func (h *HTTP) Stop() {
	h.stopOnce.Do(func() {
		close(h.sendChan)
	})
}

// Name returns the channel that the event generators will send
// events to.
func (h *HTTP) SendChannel() chan<- []byte {
	return h.sendChan
}
