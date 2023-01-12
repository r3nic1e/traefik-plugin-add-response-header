package traefik_plugin_add_response_header

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
)

var (
	_ interface {
		http.ResponseWriter
		http.Hijacker
	} = &wrappedResponseWriter{}
)

type Config struct {
	From        string `json:"from,omitempty"`
	To          string `json:"to,omitempty"`
	Regexp      string `json:"regexp,omitempty"`
	Replacement string `json:"replacement,omitempty"`
	Overwrite   bool   `json:"overwrite,omitempty"`
}

func CreateConfig() *Config {
	return &Config{
		Regexp:      "^(.*)$",
		Replacement: "$1",
	}
}

type plugin struct {
	name   string
	next   http.Handler
	config *Config
	regex  *regexp.Regexp
}

type wrappedResponseWriter struct {
	w    http.ResponseWriter
	buf  *bytes.Buffer
	code int
}

func (w *wrappedResponseWriter) Header() http.Header {
	return w.w.Header()
}

func (w *wrappedResponseWriter) Write(b []byte) (int, error) {
	return w.buf.Write(b)
}

func (w *wrappedResponseWriter) WriteHeader(code int) {
	w.code = code
}

func (w *wrappedResponseWriter) Flush() {
	w.w.WriteHeader(w.code)
	io.Copy(w.w, w.buf)
}

func (w *wrappedResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := w.w.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("%T is not an http.Hijacker", w.w)
	}

	return hijacker.Hijack()
}

func (p *plugin) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	resp := &wrappedResponseWriter{
		w:    w,
		buf:  &bytes.Buffer{},
		code: 200,
	}
	defer resp.Flush()

	p.next.ServeHTTP(resp, req)

	if !p.config.Overwrite && resp.Header().Get(p.config.To) != "" {
		return
	}

	src := req.Header.Get(p.config.From)
	if src == "" {
		return
	}

	var replacement []byte
	for _, match := range p.regex.FindAllStringSubmatchIndex(src, -1) {
		replacement = p.regex.ExpandString(
			replacement,
			p.config.Replacement,
			src,
			match,
		)
	}

	if len(replacement) > 0 {
		resp.Header().Set(p.config.To, string(replacement))
	}
}

func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config.From == "" {
		return nil, fmt.Errorf("from cannot be empty")
	}
	if config.To == "" {
		return nil, fmt.Errorf("to cannot be empty")
	}

	regex, err := regexp.Compile(config.Regexp)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regexp: %w", err)
	}

	return &plugin{
		name:   name,
		next:   next,
		config: config,
		regex:  regex,
	}, nil
}
