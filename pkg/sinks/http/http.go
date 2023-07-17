package http

import (
	"bytes"
	"net/http"
	"net/url"
	"time"

	"ctx.sh/genie/pkg/config"
	"ctx.sh/genie/pkg/sinks"
	"ctx.sh/genie/pkg/template"
)

type HTTP struct {
	url     *url.URL
	headers map[string]template.Template
	client  http.Client // change me?
	backoff time.Duration
	timeout time.Duration
	method  string
}

func FromConfig(cfg config.Configs) sinks.Sink {
	return &HTTP{}
}

func (h *HTTP) Connect() {
	h.client = http.Client{
		Timeout: h.timeout,
	}
}

func (h *HTTP) Init() {}

func (h *HTTP) Send(data []byte) {
	headers := make(map[string]string)
	for k, v := range h.headers {
		headers[k] = v.Execute()
	}
	req := http.Request{
		Method: h.method,
		URL:    h.url,
		Body:   bytes.NewBuffer(data),
	}

}
