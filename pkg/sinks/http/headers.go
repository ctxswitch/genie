package http

import (
	"stvz.io/genie/pkg/template"
)

type Headers map[string]*template.Template

func newHeaders(cfg []HttpHeaderConfig) Headers {
	headers := make(Headers)
	for _, h := range cfg {
		tmpl := template.NewTemplate()
		err := tmpl.Compile(h.Value)
		if err != nil {
			// TODO: log
			continue
		}
		headers[h.Name] = tmpl
	}
	return headers
}

func (h *Headers) Get(key string) (value any, ok bool) {
	value, ok = (*h)[key]
	return
}
