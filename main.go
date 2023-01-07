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
	if w.Header().Get("Trailer") == "" {
		w.Header().Set("Trailer", p.config.To)
	} else {
		w.Header().Add("Trailer", p.config.To)
	}

	os.Stdout.WriteString(fmt.Sprintf("ServeHTTP: w headers before - %+v", w.Header()))

	p.next.ServeHTTP(w, req)

	src := req.Header.Get(p.config.From)
	os.Stdout.WriteString(fmt.Sprintf("ServeHTTP: src header - %+v", src))

	w.Header().Set(p.config.To, src)

	os.Stdout.WriteString(fmt.Sprintf("ServeHTTP: w headers after - %+v", w.Header()))
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
