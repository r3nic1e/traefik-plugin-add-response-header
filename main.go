package traefik_plugin_add_response_header

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

type Config struct {
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
}

func CreateConfig() *Config {
	return &Config{}
}

type plugin struct {
	name   string
	next   http.Handler
	config *Config
}

type wrappedResponseWriter struct {
	w    http.ResponseWriter
	buf  *bytes.Buffer
	code int
}

func (w wrappedResponseWriter) Header() http.Header {
	return w.w.Header()
}

func (w wrappedResponseWriter) Write(b []byte) (int, error) {
	return w.buf.Write(b)
}

func (w wrappedResponseWriter) WriteHeader(code int) {
	w.code = code
}

func (w wrappedResponseWriter) Flush() {
	io.Copy(w.w, w.buf)
	w.w.WriteHeader(w.code)
}

func (p *plugin) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	resp := wrappedResponseWriter{
		w:    w,
		buf:  &bytes.Buffer{},
		code: 200,
	}

	p.next.ServeHTTP(resp, req)

	resp.Header().Set(p.config.To, req.Header.Get(p.config.From))
	resp.Flush()
}

func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config.From == "" {
		return nil, fmt.Errorf("from cannot be empty")
	}
	if config.To == "" {
		return nil, fmt.Errorf("to cannot be empty")
	}

	return &plugin{
		name:   name,
		next:   next,
		config: config,
	}, nil
}
