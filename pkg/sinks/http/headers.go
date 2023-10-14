package http

import (
	"stvz.io/genie/pkg/template"
)

// Headers is a map of names with the corresponding templates.  Headers
// can use the same templating system as events so headers can define
// dynamic values.
type Headers map[string]*template.Template

// newHeaders creates a new Headers map from a slice of HeaderConfig.
// TODO: I should be returning an error on compilation failure but I'm not
// doing that right now.
func newHeaders(cfg []HeaderConfig) Headers {
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

// Get returns the value of a header by name.
func (h *Headers) Get(key string) (value any, ok bool) {
	value, ok = (*h)[key]
	return
}
