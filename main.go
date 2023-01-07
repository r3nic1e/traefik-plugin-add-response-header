package traefik_plugin_add_response_header

import (
	"context"
	"fmt"
	"net/http"
	"os"
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

func (p *plugin) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Trailer", p.config.To)

	os.Stdout.WriteString(fmt.Sprintf("ServeHTTP: request headers before - %+v", req.Header))
	p.next.ServeHTTP(w, req)
	os.Stdout.WriteString(fmt.Sprintf("ServeHTTP: request headers after - %+v", req.Header))

	src := req.Header.Get(p.config.From)
	w.Header().Set(p.config.To, src)
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
