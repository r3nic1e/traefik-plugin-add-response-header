package traefik_plugin_add_response_header

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
}

type dummyHandler struct{}

func (dummyHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("overwrite", "true")
	w.WriteHeader(200)
}

func (s *Suite) TestMissingSourceHeader() {
	cfg := CreateConfig()
	cfg.From = "blabla"
	cfg.To = "123bla"

	h, err := New(context.Background(), dummyHandler{}, cfg, "")
	s.Require().NoError(err)

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	h.ServeHTTP(resp, req)

	s.Require().Empty(resp.Header().Get(cfg.To))
}

func (s *Suite) TestCorrectCopyWithoutOverwrite() {
	cfg := CreateConfig()
	cfg.From = "blabla"
	cfg.To = "overwrite"
	cfg.Overwrite = false

	data := "123bla321"

	h, err := New(context.Background(), dummyHandler{}, cfg, "")
	s.Require().NoError(err)

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set(cfg.From, data)
	h.ServeHTTP(resp, req)

	s.Require().NotEmpty(resp.Header().Get(cfg.To))
	s.Require().NotEqual(data, resp.Header().Get(cfg.To))
}

func (s *Suite) TestCorrectCopyWithOverwrite() {
	cfg := CreateConfig()
	cfg.From = "blabla"
	cfg.To = "overwrite"
	cfg.Overwrite = true

	data := "123bla321"

	h, err := New(context.Background(), dummyHandler{}, cfg, "")
	s.Require().NoError(err)

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set(cfg.From, data)
	h.ServeHTTP(resp, req)

	s.Require().Equal(data, resp.Header().Get(cfg.To))
}

func (s *Suite) TestRegexp() {
	cfg := CreateConfig()
	cfg.From = "blabla"
	cfg.To = "123bla"
	cfg.Regexp = "^(\\d+).*$"

	data := "123bla321"

	h, err := New(context.Background(), dummyHandler{}, cfg, "")
	s.Require().NoError(err)

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set(cfg.From, data)
	h.ServeHTTP(resp, req)

	s.Require().Equal("123", resp.Header().Get(cfg.To))
}

func (s *Suite) TestReplacement() {
	cfg := CreateConfig()
	cfg.From = "blabla"
	cfg.To = "123bla"
	cfg.Replacement = "replacement312"

	data := "123bla321"

	h, err := New(context.Background(), dummyHandler{}, cfg, "")
	s.Require().NoError(err)

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set(cfg.From, data)
	h.ServeHTTP(resp, req)

	s.Require().Equal(cfg.Replacement, resp.Header().Get(cfg.To))
}

func TestSuite(t *testing.T) {
	suite.Run(t, &Suite{})
}
