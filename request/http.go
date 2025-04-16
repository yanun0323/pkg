package helper

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

type request struct {
	ctx    context.Context
	method string
	url    string
	query  map[string]string
	header map[string]string
	body   any
	hooks  []func(*http.Request) error
}

// New creates a new request creator.
func New(method, url string) *request {
	return &request{
		method: method,
		url:    url,
		query:  make(map[string]string),
		header: make(map[string]string),
	}
}

// WithContext sets the context for the request.
func (r *request) WithContext(ctx context.Context) *request {
	r.ctx = ctx
	return r
}

// WithQueryParam sets the query parameters for the request.
func (r *request) WithQueryParam(key, value string) *request {
	r.query[key] = value
	return r
}

// WithHeader sets the header for the request.
func (r *request) WithHeader(key, value string) *request {
	r.header[key] = value
	return r
}

// WithBody sets the body for the request.
func (r *request) WithBody(p any) *request {
	r.body = p
	return r
}

// WithHook sets the hook for the request.
func (r *request) WithHook(hook func(*http.Request) error) *request {
	r.hooks = append(r.hooks, hook)
	return r
}

// Create creates a new request.
func (r *request) Create() (*http.Request, error) {
	var reader io.Reader
	if r.body != nil {
		data, err := json.Marshal(r.body)
		if err != nil {
			return nil, errors.Errorf("marshal body, err: %+v", err)
		}

		reader = bytes.NewReader(data)
	}

	if r.ctx == nil {
		r.ctx = context.Background()
	}

	req, err := http.NewRequestWithContext(r.ctx, r.method, r.url, reader)
	if err != nil {
		return nil, errors.Errorf("new request, err: %+v", err)
	}

	q := url.Values{}
	for k, v := range r.query {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	for k, v := range r.header {
		req.Header.Add(k, v)
	}

	for i, hook := range r.hooks {
		if err := hook(req); err != nil {
			return nil, errors.Errorf("execute hook(%d), err: %+v", i, err)
		}
	}

	return req, nil
}

// Send sends the request and returns the response.
//
// If proxy is provided, it will be used to send the request.
func (r *request) Send(proxy ...func(*http.Request) (*http.Response, error)) (*http.Response, error) {
	client := http.DefaultClient.Do
	if len(proxy) != 0 && proxy[0] != nil {
		client = proxy[0]
	}

	req, err := r.Create()
	if err != nil {
		return nil, err
	}

	resp, err := client(req)
	if err != nil {
		return nil, errors.Errorf("do request, err: %+v", err)
	}

	return resp, nil
}
